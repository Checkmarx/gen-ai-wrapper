package wrapper

import (
	"github.com/Checkmarx/gen-ai-wrapper/internal"
	"github.com/Checkmarx/gen-ai-wrapper/pkg/message"
	"github.com/Checkmarx/gen-ai-wrapper/pkg/models"
)

// LitellmWrapper provides a simple wrapper for litellm AI proxy service
type LitellmWrapper interface {
	Call(cxAuth string, metaData *message.MetaData, request *internal.ChatCompletionRequest) (*internal.ChatCompletionResponse, error)
	SetupCall(messages []message.Message)
	Close() error
}

// LitellmWrapperImpl implements the LitellmWrapper interface
type LitellmWrapperImpl struct {
	wrapper internal.Wrapper
	model   string
}

// NewLitellmWrapper creates a new litellm wrapper
func NewLitellmWrapper(endPoint, apiKey, model string) (LitellmWrapper, error) {
	if model == "" {
		model = models.DefaultModel
	}

	wrapper := internal.NewLitellmWrapper(endPoint, apiKey, 0)

	return &LitellmWrapperImpl{
		wrapper: wrapper,
		model:   model,
	}, nil
}

// SetupCall sets up the wrapper with initial messages
func (w *LitellmWrapperImpl) SetupCall(messages []message.Message) {
	w.wrapper.SetupCall(messages)
}

// Call makes a request to the litellm service
func (w *LitellmWrapperImpl) Call(cxAuth string, metaData *message.MetaData, request *internal.ChatCompletionRequest) (*internal.ChatCompletionResponse, error) {
	// Set the model if not already set
	if request.Model == "" {
		request.Model = w.model
	}

	return w.wrapper.Call(cxAuth, metaData, request)
}

// Close closes the wrapper
func (w *LitellmWrapperImpl) Close() error {
	return w.wrapper.Close()
}
