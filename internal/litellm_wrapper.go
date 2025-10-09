package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Checkmarx/gen-ai-wrapper/pkg/message"
)

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

// Call makes a request to the litellm AI proxy service
func (w *LitellmWrapper) Call(cxAuth string, metaData *message.MetaData, request *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	// Prepare the request
	req, err := w.prepareRequest(cxAuth, metaData, request)
	if err != nil {
		return nil, err
	}

	// Make the HTTP request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Handle the response
	return w.handleResponse(resp)
}

// prepareRequest creates the HTTP request
func (w *LitellmWrapper) prepareRequest(cxAuth string, metaData *message.MetaData, requestBody *ChatCompletionRequest) (*http.Request, error) {
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, w.endPoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cxAuth))

	// Set required headers for litellm service
	req.Header.Set("X-Request-ID", metaData.RequestID)
	req.Header.Set("X-Tenant-ID", metaData.TenantID)
	req.Header.Set("User-Agent", metaData.UserAgent)
	req.Header.Set("X-Feature", metaData.Feature)

	return req, nil
}

// handleResponse processes the HTTP response
func (w *LitellmWrapper) handleResponse(resp *http.Response) (*ChatCompletionResponse, error) {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Handle successful response
	if resp.StatusCode == http.StatusOK {
		var responseBody = new(ChatCompletionResponse)
		err = json.Unmarshal(bodyBytes, responseBody)
		if err != nil {
			return nil, err
		}
		return responseBody, nil
	}

	// Handle error responses
	var errorResponse = new(ErrorResponse)
	err = json.Unmarshal(bodyBytes, errorResponse)
	if err != nil {
		// If we can't parse the error response, return a generic error
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Return the parsed error
	return nil, fromResponse(resp.StatusCode, errorResponse)
}

// Close closes the wrapper (no-op for HTTP client)
func (w *LitellmWrapper) Close() error {
	return nil
}
