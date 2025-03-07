package dnhtclient

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/andybalholm/brotli"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	client *http.Client
}

func NewClient() *Client {
	customTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &Client{
		client: &http.Client{Transport: customTransport},
	}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

func ExecuteRequest(client HTTPClient, method, url string, headers http.Header, body string) (string, error) {
	req, err := http.NewRequest(method, url, bytes.NewBufferString(body))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}
	req.Header = headers

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	respBody := new(strings.Builder)
	if resp.Header.Get("Content-Encoding") == "br" {
		br := brotli.NewReader(resp.Body)
		_, err = bufio.NewReader(br).WriteTo(respBody)
	} else {
		_, err = bufio.NewReader(resp.Body).WriteTo(respBody)
	}
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	return strconv.Itoa(respBody.Len()), nil
}
