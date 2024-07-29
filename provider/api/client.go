package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	APIUrl    string
	APIKey    string
	UserEmail string
}

func NewClient(apiURL, apiKey, userEmail string) (*Client, error) {
	return &Client{
		APIUrl:    apiURL,
		APIKey:    apiKey,
		UserEmail: userEmail,
	}, nil
}

func (c *Client) apiRequest(method, endpoint string, payload interface{}) (*http.Response, error) {
	var body io.Reader
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("error marshaling JSON: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	url := fmt.Sprintf("%s%s", c.APIUrl, endpoint)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("x-mcs-user", c.UserEmail)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("error response from API: %s - %s", resp.Status, string(bodyBytes))
	}

	return resp, nil
}
