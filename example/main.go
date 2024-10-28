package main

import (
	"flag"
	"fmt"
	"github.com/Checkmarx/gen-ai-wrapper/pkg/connector"
	"github.com/Checkmarx/gen-ai-wrapper/pkg/message"
	"github.com/Checkmarx/gen-ai-wrapper/pkg/role"
	"github.com/Checkmarx/gen-ai-wrapper/pkg/wrapper"
	"github.com/google/uuid"
	"os"
	"strings"
)

const usage = `
Chat with GPT

Usage: chat [-s <system-prompt>] -u <user-prompt> [options]  
   or: chat -id <conversation-id> -u <user-prompt> [options]

Options
   -s, --system <system-prompt>  system prompt string
   -u, --user <user-prompt>      user prompt string
   -id <conversation-id>         chat conversation ID
   -ai <ai-server>               AI server to use. Options: {OpenAI (default), CxOne}
   -f, --full-response           return full response from GPT
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

	err = CallAIandPrintResponse(aiServer, systemPrompt, userPrompt, id, fullResponse)
	if err != nil {
		fmt.Printf("Error '%v' calling GPT\n", err)
		os.Exit(1)
	}

}

func CallAIandPrintResponse(aiServer, systemPrompt, userPrompt string, chatId uuid.UUID, fullResponse bool) error {

	aiKey, err := getAIAccessKey(aiServer)
	if err != nil {
		return err
	}
	aiEndpoint, err := getAIEndpoint(aiServer)
	if err != nil {
		return err
	}

	statefulWrapper, err := wrapper.NewStatefulWrapperNew(
		connector.NewFileSystemConnector(""), aiEndpoint, aiKey, "gpt-4o", 4, 0)
	if err != nil {
		return fmt.Errorf("error creating GPT client: %v", err)
	}

	var newMessages []message.Message
	newMessages = append(newMessages, message.Message{
		Role:    role.System,
		Content: systemPrompt,
	})
	newMessages = append(newMessages, message.Message{
		Role:    role.User,
		Content: userPrompt,
	})

	if fullResponse {
		response, err := statefulWrapper.SecureCallReturningFullResponse("", nil, chatId, newMessages)
		if err != nil {
			return fmt.Errorf("error calling GPT: %v", err)
		}
		fmt.Printf("%+v\n", response)
	} else {
		response, err := statefulWrapper.Call(chatId, newMessages)
		if err != nil {
			return fmt.Errorf("error calling GPT: %v", err)
		}
		fmt.Println(getMessageContents(response))
	}
	return nil
}

func getAIAccessKey(aiServer string) (string, error) {
	if strings.EqualFold(aiServer, "OpenAI") {
		accessKey, err := GetOpenAIAccessKey()
		if err != nil {
			return "", fmt.Errorf("error getting OpenAI API key: %v", err)
		}
		return accessKey, nil
	}
	if strings.EqualFold(aiServer, "CxOne") {
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
	if strings.EqualFold(aiServer, "CxOne") {
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
