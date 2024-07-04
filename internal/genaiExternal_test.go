package internal

import (
	"bytes"
	"encoding/json"
	"github.com/Checkmarx/gen-ai-wrapper/pkg/message"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strconv"
	"testing"
)

func TestHandleGptResponse(t *testing.T) {
	wrapper := WrapperImpl{
		apiKey:   "test",
		endPoint: "testServer.URL",
		dropLen:  4}
	resp := &http.Response{Body: io.NopCloser(bytes.NewReader([]byte(RespBody)))}
	resp.StatusCode = http.StatusOK
	_, err := wrapper.handleGptResponse("test", &message.MetaData{}, &ChatCompletionRequest{}, resp)
	if err != nil {
		t.Fatal(err)
	}
}

func TestHandleGptResponseNegativeOpenAi(t *testing.T) {
	wrapper := WrapperImpl{
		apiKey:   "test",
		endPoint: "testServer.URL",
		dropLen:  4}
	gptError := GptError{
		Message: "test",
		Type:    "test",
		Param:   "test",
		Code:    429,
	}
	errRes := ErrorResponse{
		Error: gptError,
	}
	errResB, _ := json.Marshal(errRes)
	resp := &http.Response{Body: io.NopCloser(bytes.NewReader(errResB))}
	resp.StatusCode = http.StatusTooManyRequests
	res, err := wrapper.handleGptResponse("test", nil, &ChatCompletionRequest{}, resp)
	if res != nil {
		t.Fatal("Expected nil response")
	}
	if err == nil {
		t.Fatal("Expected error")
	}
	assert.Contains(t, err.Error(), "test")
	assert.Contains(t, err.Error(), strconv.Itoa(http.StatusTooManyRequests))
}

func TestHandleGptResponseNegativeAzureExternal(t *testing.T) {
	wrapper := WrapperImpl{
		apiKey:   "test",
		endPoint: "testServer.URL",
		dropLen:  4}
	gptError := GptError{
		Message: "test",
		Type:    "test",
		Param:   "test",
		Code:    429,
	}
	errRes := ErrorResponse{
		Error: gptError,
	}
	errResB, _ := json.Marshal(errRes)
	resp := &http.Response{Body: io.NopCloser(bytes.NewReader(errResB))}
	resp.Header = make(http.Header)
	resp.Header.Set("X-Gen-Ai-ErrorCode", strconv.Itoa(http.StatusTooManyRequests))
	resp.StatusCode = http.StatusFailedDependency
	res, err := wrapper.handleGptResponse("test", &message.MetaData{
		RequestID: "test",
		TenantID:  "test",
		UserAgent: "test",
		Feature:   "test",
	}, &ChatCompletionRequest{}, resp)
	if res != nil {
		t.Fatal("Expected nil response")
	}
	if err == nil {
		t.Fatal("Expected error")
	}
	assert.Contains(t, err.Error(), "test")
	assert.Contains(t, err.Error(), strconv.Itoa(http.StatusTooManyRequests))
	assert.NotContains(t, err.Error(), strconv.Itoa(http.StatusFailedDependency))
}

func TestHandleGptResponseNegativeAzureInternal(t *testing.T) {
	wrapper := WrapperImpl{
		apiKey:   "test",
		endPoint: "testServer.URL",
		dropLen:  4}
	gptError := GptError{
		Message: "test",
		Type:    "test",
		Param:   "test",
		Code:    429,
	}
	errRes := ErrorResponse{
		Error: gptError,
	}
	errResB, _ := json.Marshal(errRes)
	resp := &http.Response{Body: io.NopCloser(bytes.NewReader(errResB))}
	resp.Header = make(http.Header)
	resp.Header.Set("X-Gen-Ai-ErrorCode", strconv.Itoa(0))
	resp.StatusCode = http.StatusInternalServerError
	res, err := wrapper.handleGptResponse("test", &message.MetaData{
		RequestID: "test",
		TenantID:  "test",
		UserAgent: "test",
		Feature:   "test",
	}, &ChatCompletionRequest{}, resp)
	if res != nil {
		t.Fatal("Expected nil response")
	}
	if err == nil {
		t.Fatal("Expected error")
	}
	assert.NotContains(t, err.Error(), "test")
	assert.Contains(t, err.Error(), strconv.Itoa(0))
	assert.Contains(t, err.Error(), strconv.Itoa(http.StatusInternalServerError))
}

//nolint:lll
const (
	RespBody = `{
"id":"chatcmpl-8v26ZSbJu0DEIbeyLVG3NT8m7imGl",
"object":"chat.completion","created":1708603479,
"model":"gpt-4","choices":[{
"finish_reason":"stop","index":0,"message":
{"role":"assistant","content":"The impact of the 'ALB Deletion Protection Disabled' issue is potentially significant, depending on the criticality of the AWS Application Load Balancer (ALB) to your applications or services. Since 'enable_deletion_protection' is currently set to 'false', it means deletion protection is not enabled on your ALB. This allows anyone with the appropriate permissions to delete your ALB. If this load balancer is deleted, either accidentally or maliciously, it can lead to service disruptions as the traffic meant for your applications will not be distributed. Particularly in a production environment, such a service disruption can lead to loss of availability or potentially even loss of business or revenue, especially if the recovery and restoration time is lengthy. This is why it's recommended to set the 'enable_deletion_protection' field to 'true', especially when dealing with ALBs in production environments. "}}],
"usage":{"prompt_tokens":849,"completion_tokens":181,"total_tokens":1030}}
`
)
