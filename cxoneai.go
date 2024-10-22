package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

var (
	accessToken    string
	expirationTime time.Time
	tokenMutex     sync.Mutex
)

const (
	cx1OAuthURLEnv     = "CXONE_OAUTH_URL"
	cx1ClientIDEnv     = "CXONE_CLIENT_ID"
	cx1ClientSecretEnv = "CXONE_CLIENT_SECRET"
	cx1AICompletionUrl = "CXONE_AI_URL"
)

type OAuthTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func GetCxOneAIAccessKey() (string, error) {
	if accessToken == "" {
		if err := getOAuthAccessToken(); err != nil {
			return "", err
		}
		StartTokenRefresher()
	}
	return accessToken, nil

}

func GetCxOneAIEndpoint() (string, error) {
	endpoint, err := GetEnvKeyValue(cx1AICompletionUrl)
	if err != nil {
		return "", err
	}

	return endpoint, nil
}

func getOAuthAccessToken() error {
	openIDURL, err := GetEnvKeyValue(cx1OAuthURLEnv)
	if err != nil {
		return err
	}
	clientID, err := GetEnvKeyValue(cx1ClientIDEnv)
	if err != nil {
		return err
	}
	clientSecret, err := GetEnvKeyValue(cx1ClientSecretEnv)
	if err != nil {
		return err
	}

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	req, err := http.NewRequest("POST", openIDURL, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to get access token")
	}

	var tokenResponse OAuthTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return err
	}

	tokenMutex.Lock()
	expirationTime = time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second)
	accessToken = tokenResponse.AccessToken
	tokenMutex.Unlock()
	return nil
}

func StartTokenRefresher() {
	go func() {
		for {
			time.Sleep(time.Until(expirationTime.Add(-1 * time.Minute)))
			err := getOAuthAccessToken()
			if err != nil {
				fmt.Printf("Failed to refresh CxOne OAuth token: %v\n", err)
				continue
			}
			fmt.Println("Refreshed CxOne OAuth token")
		}
	}()
}
