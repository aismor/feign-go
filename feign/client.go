package feign

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	BaseURL    string
}

func NewClient(baseURL string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		BaseURL:    baseURL,
	}
}

func NewClientWithToken(baseURL string, tokenProvider TokenProvider) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout:   10 * time.Second,
			Transport: &TokenTransport{TokenProvider: tokenProvider},
		},
		BaseURL: baseURL,
	}
}

func (c *Client) Get(path string, respBody interface{}) error {
	fullURL := c.BaseURL + path
	resp, err := c.httpClient.Get(fullURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return decodeResponse(resp.Body, respBody)
}

func (c *Client) Post(path string, reqBody interface{}, respBody interface{}) error {
	return c.doRequestWithBody(http.MethodPost, path, reqBody, respBody)
}

func (c *Client) Put(path string, reqBody interface{}, respBody interface{}) error {
	return c.doRequestWithBody(http.MethodPut, path, reqBody, respBody)
}

func (c *Client) Delete(path string, respBody interface{}) error {
	req, err := http.NewRequest(http.MethodDelete, c.BaseURL+path, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return decodeResponse(resp.Body, respBody)
}

func (c *Client) doRequestWithBody(method string, path string, reqBody interface{}, respBody interface{}) error {
	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, c.BaseURL+path, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return decodeResponse(resp.Body, respBody)
}
