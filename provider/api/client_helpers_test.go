package api

func testClient(apiUrl string) *Client {
	return &Client{
		APIUrl:    apiUrl,
		APIKey:    "dummy-key",
		UserEmail: "user@example.com",
	}
}
