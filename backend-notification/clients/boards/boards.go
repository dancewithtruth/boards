package boards

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type BoardsClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(baseURL string) *BoardsClient {
	return &BoardsClient{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
}

func (c *BoardsClient) CreateEmailVerification(input CreateEmailVerificationInput) (*http.Response, error) {
	return c.makeRequest("POST", "/users/email-verifications", input)
}

func (c *BoardsClient) GetInvite(inviteID string) (*http.Response, error) {
	return c.makeRequest("GET", "/invites/"+inviteID, nil)
}

// makeRequest sends an HTTP request and returns the response.
func (c *BoardsClient) makeRequest(method, endpoint string, payload interface{}) (*http.Response, error) {
	url := c.BaseURL + endpoint

	var requestBody []byte
	if payload != nil {
		requestBody, _ = json.Marshal(payload)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
