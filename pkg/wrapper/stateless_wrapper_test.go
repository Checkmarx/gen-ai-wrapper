package wrapper

import (
	"fmt"
	"testing"

	"github.com/Checkmarx/gen-ai-wrapper/pkg/message"
	"github.com/Checkmarx/gen-ai-wrapper/pkg/models"
	"github.com/Checkmarx/gen-ai-wrapper/pkg/role"
)

func TestCallGPT(t *testing.T) {
	var history []message.Message
	var response []message.Message
	wrapper, err := NewStatelessWrapper(OpenAiEndPoint, apikey, models.GPT3Dot5Turbo, 4, 0)
	if err != nil {
		t.Fatal(err)
	}
	for _, q := range userQuestions {
		t.Log(q)
		var newMessages []message.Message
		newMessages = append(newMessages, message.Message{
			Role:    role.System,
			Content: systemInput,
		})
		newMessages = append(newMessages, message.Message{
			Role:    role.Assistant,
			Content: assistantInput,
		})
		newMessages = append(newMessages, message.Message{
			Role:    role.User,
			Content: fmt.Sprintf(userInput, q),
		})

		response, err = wrapper.Call(history, newMessages)
		if err != nil {
			t.Fatal(err)
		}
		if len(response) != 1 {
			t.Fatalf("Got multiple choices\n%v\n", response)
		}

		history = append(history, newMessages...)
		history = append(history, response[0])
	}
	for _, m := range history {
		t.Logf("%s\n\n%s\n\n", m.Role, m.Content)
	}
}

func TestCallEmptyApiKey(t *testing.T) {
	wrapper, err := NewStatelessWrapper(OpenAiEndPoint, apikey, models.GPT3Dot5Turbo, 4, 0)
	if err != nil {
		t.Fatal(err)
	}
	q := userQuestions[0]
	_, err = wrapper.Call(nil, []message.Message{{
		Role:    role.User,
		Content: fmt.Sprintf(userInput, q),
	}})
	if err == nil {
		t.Fatal("Call succeeded without API key")
	}
}
