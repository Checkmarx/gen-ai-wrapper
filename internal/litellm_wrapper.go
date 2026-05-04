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

var httpClient = &http.Client{Timeout: 120 * time.Second}

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
func validateEndpoint(endpoint string) error {
	parsed, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("invalid endpoint URL: %w", err)
	}
	if parsed.Scheme != "https" {
		return fmt.Errorf("endpoint must use https scheme, got: %q", parsed.Scheme)
	}
	host := parsed.Hostname()
	if host == "" {
		return fmt.Errorf("endpoint URL has no host")
	}
	ip := net.ParseIP(host)
	if ip != nil {
		if ip.IsLoopback() || ip.IsPrivate() || ip.IsUnspecified() || ip.IsLinkLocalUnicast() {
			return fmt.Errorf("endpoint host resolves to a disallowed address: %s", host)
		}
	} else {
		lower := strings.ToLower(host)
		if lower == "localhost" || strings.HasSuffix(lower, ".local") || strings.HasSuffix(lower, ".internal") {
			return fmt.Errorf("endpoint host is a disallowed internal address: %s", host)
		}
	}
	return nil
}

// Call makes a request to the litellm AI proxy service
func (w *LitellmWrapper) Call(cxAuth string, metaData *message.MetaData, request *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	if err := validateEndpoint(w.endPoint); err != nil {
		return nil, err
	}

	req, err := w.prepareRequest(cxAuth, metaData, request)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return w.handleResponse(resp)
}

// prepareRequest creates the HTTP request with an explicit timeout context
func (w *LitellmWrapper) prepareRequest(cxAuth string, metaData *message.MetaData, requestBody *ChatCompletionRequest) (*http.Request, error) {
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	_ = cancel // caller's defer resp.Body.Close triggers GC; context is bounded by the client timeout

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, w.endPoint, bytes.NewBuffer(jsonData))
	if err != nil {
		cancel()
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cxAuth))
	req.Header.Set("X-Request-ID", metaData.RequestID)
	req.Header.Set("X-Tenant-ID", metaData.TenantID)
	req.Header.Set("User-Agent", metaData.UserAgent)
	req.Header.Set("X-Feature", metaData.Feature)

	return req, nil
}

// handleResponse processes the HTTP response
func (w *LitellmWrapper) handleResponse(resp *http.Response) (*ChatCompletionResponse, error) {
	bodyBytes, err := io.ReadAll(io.LimitReader(resp.Body, 10*1024*1024))
	if err != nil {
		return nil, err
	}

	if len(bodyBytes) == 0 {
		return nil, fmt.Errorf("HTTP %d: empty response body", resp.StatusCode)
	}

	if resp.StatusCode == http.StatusOK {
		if !json.Valid(bodyBytes) {
			return nil, fmt.Errorf("HTTP %d: response body is not valid JSON", resp.StatusCode)
		}
		var responseBody ChatCompletionResponse
		if err := json.Unmarshal(bodyBytes, &responseBody); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}
		return &responseBody, nil
	}

	if !json.Valid(bodyBytes) {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(bodyBytes))
	}
	var errorResponse ErrorResponse
	if err := json.Unmarshal(bodyBytes, &errorResponse); err != nil {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil, fromResponse(resp.StatusCode, &errorResponse)
}

// Close closes the wrapper (no-op for HTTP client)
func (w *LitellmWrapper) Close() error {
	return nil
}
