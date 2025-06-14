package api

func MockClient(apiUrl string) *Client {
	return &Client{
		APIUrl:    apiUrl,
		AuthType:  "apiKey",
		APIKey:    "dummy-key",
		UserEmail: "user@example.com",
	}
}
