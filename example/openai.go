package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	OpenAICompletionUrl = "OPENAI_COMPLETION_URL"
	OpenAIEnginesUrl    = "OPENAI_ENGINES_URL"
	OpenAIApiKey        = "OPENAI_API_KEY"
)

func GetOpenAIAccessKey(model string) (string, error) {
	apiKey, err := GetEnvKeyValue(OpenAIApiKey)
	if err != nil {
		return "", err
	}

	if err := isValidAIApiKey(apiKey, model); err != nil {
		return "", err
	}
	return apiKey, nil
}

func GetOpenAIEndpoint() (string, error) {
	endpoint, err := GetEnvKeyValue(OpenAICompletionUrl)
	if err != nil {
		return "", err
	}

	return endpoint, nil
}

type EngineInfo struct {
	ID            string `json:"id"`
	Object        string `json:"object"`
	Created       int    `json:"created"`
	Description   string `json:"description"`
	MaxTokens     int    `json:"max_tokens"`
	Name          string `json:"name"`
	ReadyAvail    bool   `json:"ready_availability"`
	Owner         string `json:"owner"`
	Permissions   string `json:"permissions"`
	Plan          string `json:"plan"`
	PricePerToken string `json:"price_per_token"`
}

type ErrorMessage struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Param   string `json:"param"`
		Code    string `json:"code"`
	} `json:"error"`
}

func isValidAIApiKey(apiKey, model string) error {
	enginesUrl, err := GetEnvKeyValue(OpenAIEnginesUrl)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/%s", enginesUrl, model)
	request, _ := http.NewRequest("GET", url, nil)

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(request)

	if err != nil { // response error
		return err
	} else if resp.StatusCode != http.StatusOK { // invalid key
		var errMsg ErrorMessage
		body, _ := io.ReadAll(resp.Body)

		err := json.Unmarshal([]byte(body), &errMsg)
		if err == nil {
			return fmt.Errorf("OpenAI API key does not support '%s': '%s'", model, errMsg.Error.Message)
		} else {
			return err
		}
	} else { // response is OK
		var engineInfo EngineInfo
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		} else if err := json.Unmarshal(body, &engineInfo); err != nil {
			return err
		} else if strings.EqualFold(engineInfo.ID, model) {
			return nil
		} else {
			return fmt.Errorf("model '%s' is inaccessible with the given API Key", model)
		}
	}
}
