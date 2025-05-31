package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Package utils provides a simple HTTP client for making GET and POST requests.
type Client struct {
	httpClient *http.Client
	BaseURL    string
	Headers    map[string]string
	Timeout    time.Duration
}

// New creates a new HTTP client with the given base URL, timeout, and headers.
func NewClient(baseURL string, timeout time.Duration, headers map[string]string) *Client {
	return &Client{
		BaseURL: baseURL,
		Headers: headers,
		Timeout: timeout,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// Get performs a GET request to the specified path with optional query parameters.
func (c *Client) Get(ctx context.Context, path string, query map[string]string) ([]byte, error) {
	fmt.Println("Request URL:", c.BaseURL+path)
	reqURL := c.BaseURL + path
	if len(query) > 0 {
		q := url.Values{}
		for k, v := range query {
			q.Add(k, v)
		}
		reqURL += "?" + q.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}

	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return readResponse(resp)
}

// Post performs a POST request to the specified path with the given body.
func (c *Client) Post(ctx context.Context, path string, body any) ([]byte, error) {
	reqURL := c.BaseURL + path

	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}

	c.setHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return readResponse(resp)
}

// setHeaders
func (c *Client) setHeaders(req *http.Request) {
	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}
}

// readResponse
func readResponse(resp *http.Response) ([]byte, error) {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP error: %d - %s", resp.StatusCode, string(body))
	}
	return io.ReadAll(resp.Body)
}
