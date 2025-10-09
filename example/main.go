package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Checkmarx/gen-ai-wrapper/internal"
	"github.com/Checkmarx/gen-ai-wrapper/pkg/message"
	"github.com/Checkmarx/gen-ai-wrapper/pkg/models"
	"github.com/Checkmarx/gen-ai-wrapper/pkg/role"
	"github.com/Checkmarx/gen-ai-wrapper/pkg/wrapper"
	"github.com/google/uuid"
)

const usage = `
Chat with GPT

Usage: chat [-s <system-prompt>] -u <user-prompt> [options]  
   or: chat -id <conversation-id> -u <user-prompt> [options]

Options
   -s, --system <system-prompt>  system (or developer) prompt string
   -u, --user <user-prompt>      user prompt string
   -id <conversation-id>         chat conversation ID
   -ai <ai-server>               AI server to use. Options: {OpenAI (default), CxOne, LiteLLM}
   -m, --model <model>           model to use. Options: {gpt-4o (default), gpt-4, o1, o1-mini, claude-3-5-sonnet-20241022, ...}
   -h, --help                    show help
`

func main() {
	var help bool
	flag.BoolVar(&help, "help", false, "")
	flag.BoolVar(&help, "h", false, "")
	var systemPrompt string
	flag.StringVar(&systemPrompt, "system", "", "")
	flag.StringVar(&systemPrompt, "s", "", "")
	var userPrompt string
	flag.StringVar(&userPrompt, "user", "", "")
	flag.StringVar(&userPrompt, "u", "", "")

	var conversationId string
	flag.StringVar(&conversationId, "id", "", "")
	var aiServer string
	flag.StringVar(&aiServer, "ai", "OpenAI", "")
	var model string
	flag.StringVar(&model, "model", "gpt-4o", "")
	flag.StringVar(&model, "m", "gpt-4o", "")

	var fullResponse bool
	flag.BoolVar(&fullResponse, "full-response", false, "")
	flag.BoolVar(&fullResponse, "f", false, "")

	flag.Usage = func() {
		fmt.Print(usage)
		os.Exit(1)
	}

	flag.Parse()
	aiServer = strings.ToLower(aiServer)
	if help {
		printHelp()
	}

	if userPrompt == "" {
		fmt.Println("user prompt is required")
		printHelp()
	}
	if conversationId != "" && systemPrompt != "" {
		fmt.Println("system prompt cannot be specified with a conversation ID")
		printHelp()
	}

	if conversationId == "" {
		conversationId = uuid.New().String()
	}

	id, err := uuid.Parse(conversationId)
	if err != nil {
		fmt.Printf("Error parsing conversation ID: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s:\n", conversationId)

	err = CallAIandPrintResponse(aiServer, model, systemPrompt, userPrompt, id, fullResponse)
	if err != nil {
		fmt.Printf("Error '%v' calling AI\n", err)
		os.Exit(1)
	}

}

func CallAIandPrintResponse(aiServer, model, systemPrompt, userPrompt string, chatId uuid.UUID, fullResponse bool) error {

	aiKey, err := getAIAccessKey(aiServer, model)
	if err != nil {
		return err
	}
	aiEndpoint, err := getAIEndpoint(aiServer)
	if err != nil {
		return err
	}

	var litellmWrapper wrapper.LitellmWrapper

	// Use litellm wrapper for litellm server
	if strings.EqualFold(aiServer, "LiteLLM") {
		litellmWrapper, err = wrapper.NewLitellmWrapper(aiEndpoint, aiKey, model)
	} else {
		// For other servers, we'll need to implement or use existing wrappers
		return fmt.Errorf("unsupported AI server: %s", aiServer)
	}

	if err != nil {
		return fmt.Errorf("error creating '%s' AI client: %v", aiServer, err)
	}

	newMessages := GetMessages(model, systemPrompt, userPrompt)

	// Create proper metadata for the request
	metaData := &message.MetaData{
		RequestID: "example-request-" + chatId.String(),
		TenantID:  "default-tenant",
		UserAgent: "gen-ai-wrapper-example",
		Feature:   "chat-completion",
	}

	// Create the request
	request := &internal.ChatCompletionRequest{
		Model:    model,
		Messages: newMessages,
	}

	// Make the call
	response, err := litellmWrapper.Call(aiKey, metaData, request)
	if err != nil {
		return fmt.Errorf("error calling litellm: %v", err)
	}

	if fullResponse {
		fmt.Printf("%+v\n", response)
	} else {
		fmt.Println(response.Choices[0].Message.Content)
	}
	return nil
}

func GetMessages(model, systemPrompt, userPrompt string) []message.Message {
	var newMessages []message.Message
	if !strings.HasPrefix(model, "o1-") {
		r := role.System
		if model == models.O1 {
			r = role.Developer
		}
		newMessages = append(newMessages, message.Message{
			Role:    r,
			Content: systemPrompt,
		})
	} else {
		userPrompt = systemPrompt + "\n" + userPrompt
	}
	newMessages = append(newMessages, message.Message{
		Role:    role.User,
		Content: userPrompt,
	})
	return newMessages
}

func getAIAccessKey(aiServer, model string) (string, error) {
	if strings.EqualFold(aiServer, "OpenAI") {
		accessKey, err := GetOpenAIAccessKey(model)
		if err != nil {
			return "", fmt.Errorf("error getting OpenAI API key: %v", err)
		}
		return accessKey, nil
	}
	if strings.EqualFold(aiServer, "CxOne") || strings.EqualFold(aiServer, "LiteLLM") {
		accessKey, err := GetCxOneAIAccessKey()
		if err != nil {
			return "", fmt.Errorf("error getting CxOne AI API key: %v", err)
		}
		return accessKey, nil
	}
	return "", fmt.Errorf("unknown AI server: %s", aiServer)
}

func getAIEndpoint(aiServer string) (string, error) {
	if strings.EqualFold(aiServer, "OpenAI") {
		aiEndpoint, err := GetOpenAIEndpoint()
		if err != nil {
			return "", fmt.Errorf("error getting OpenAI endpoint: %v", err)
		}
		return aiEndpoint, nil
	}
	if strings.EqualFold(aiServer, "CxOne") || strings.EqualFold(aiServer, "LiteLLM") {
		aiEndpoint, err := GetCxOneAIEndpoint()
		if err != nil {
			return "", fmt.Errorf("error getting CxOne AI endpoint: %v", err)
		}
		return aiEndpoint, nil
	}
	return "", fmt.Errorf("unknown AI server: %s", aiServer)
}

func getMessageContents(response []message.Message) string {
	var responseContent []string
	for _, r := range response {
		responseContent = append(responseContent, r.Content)
	}
	return strings.Join(responseContent, "\n")
}

func printHelp() {
	flag.Usage()
	os.Exit(0)
}
