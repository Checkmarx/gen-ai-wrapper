package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Checkmarx/gen-ai-wrapper/pkg/message"
)

const httpTimeout = 120 * time.Second

var httpClient = &http.Client{Timeout: httpTimeout}

// LitellmWrapper implements the Wrapper interface for litellm AI proxy service
type LitellmWrapper struct {
	endPoint string
	apiKey   string
}

// NewLitellmWrapper creates a new litellm wrapper instance
func NewLitellmWrapper(endPoint, apiKey string) Wrapper {
	return &LitellmWrapper{
		endPoint: endPoint,
		apiKey:   apiKey,
	}
}

// SetupCall sets up the wrapper with initial messages (no-op for litellm)
func (w *LitellmWrapper) SetupCall(messages []message.Message) {
	// No setup needed for litellm
}

// validateEndpoint checks the endpoint URL to prevent SSRF attacks.
// Only https scheme is allowed and private/loopback hosts are rejected.
func validateEndpoint(endpoint string) (*url.URL, error) {
	parsed, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid endpoint URL: %w", err)
	}
	if parsed.Scheme != "https" {
		return nil, fmt.Errorf("endpoint must use https scheme, got: %q", parsed.Scheme)
	}
	host := parsed.Hostname()
	if host == "" {
		return nil, fmt.Errorf("endpoint URL has no host")
	}
	ip := net.ParseIP(host)
	if ip != nil {
		if ip.IsLoopback() || ip.IsPrivate() || ip.IsUnspecified() || ip.IsLinkLocalUnicast() {
			return nil, fmt.Errorf("endpoint host resolves to a disallowed address: %s", host)
		}
	} else {
		lower := strings.ToLower(host)
		if lower == "localhost" || strings.HasSuffix(lower, ".local") || strings.HasSuffix(lower, ".internal") {
			return nil, fmt.Errorf("endpoint host is a disallowed internal address: %s", host)
		}
	}
	return parsed, nil
}

// Call makes a request to the litellm AI proxy service
func (w *LitellmWrapper) Call(cxAuth string, metaData *message.MetaData, request *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	// Validate and parse the endpoint URL before use (SSRF prevention)
	parsedURL, err := validateEndpoint(w.endPoint)
	if err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	// Use an explicit timeout context so the deadline is visible at the call site
	ctx, cancel := context.WithTimeout(context.Background(), httpTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, parsedURL.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cxAuth))
	req.Header.Set("X-Request-ID", metaData.RequestID)
	req.Header.Set("X-Tenant-ID", metaData.TenantID)
	req.Header.Set("User-Agent", metaData.UserAgent)
	req.Header.Set("X-Feature", metaData.Feature)

	client := &http.Client{Timeout: httpTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Validate HTTP status code before processing the body
	if resp.StatusCode != http.StatusOK {
		return nil, w.handleErrorResponse(resp)
	}

	return w.handleSuccessResponse(resp)
}

// handleSuccessResponse decodes a 200 OK response
func (w *LitellmWrapper) handleSuccessResponse(resp *http.Response) (*ChatCompletionResponse, error) {
	bodyBytes, err := io.ReadAll(io.LimitReader(resp.Body, 10*1024*1024))
	if err != nil {
		return nil, err
	}
	if len(bodyBytes) == 0 {
		return nil, fmt.Errorf("HTTP %d: empty response body", resp.StatusCode)
	}
	if !json.Valid(bodyBytes) {
		return nil, fmt.Errorf("HTTP %d: response body is not valid JSON", resp.StatusCode)
	}
	var responseBody ChatCompletionResponse
	if err := json.Unmarshal(bodyBytes, &responseBody); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &responseBody, nil
}

// handleErrorResponse decodes a non-200 response into an error
func (w *LitellmWrapper) handleErrorResponse(resp *http.Response) error {
	bodyBytes, err := io.ReadAll(io.LimitReader(resp.Body, 10*1024*1024))
	if err != nil {
		return fmt.Errorf("HTTP %d: failed to read error body: %w", resp.StatusCode, err)
	}
	if len(bodyBytes) == 0 {
		return fmt.Errorf("HTTP %d: empty error response body", resp.StatusCode)
	}
	if !json.Valid(bodyBytes) {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(bodyBytes))
	}
	var errorResponse ErrorResponse
	if err := json.Unmarshal(bodyBytes, &errorResponse); err != nil {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(bodyBytes))
	}
	return fromResponse(resp.StatusCode, &errorResponse)
}

// Close closes the wrapper (no-op for HTTP client)
func (w *LitellmWrapper) Close() error {
	return nil
}
