package switchbotapi

import "net/http"

type Client struct {
	httpClient *http.Client
	token string
	secret string
}

func NewClient(httpClient *http.Client, token, secret string) *Client {
	return &Client{httpClient: httpClient, token: token, secret: secret}
}
