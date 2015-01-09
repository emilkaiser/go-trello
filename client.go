package client

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Interface describes the Trello API Client, used for mocking purposes
type Interface interface {
	Get(uri string, optionalParams ...map[string]string) ([]byte, error)
	Post(uri string, optionalParams ...map[string]string) ([]byte, error)
	Put(uri string, optionalParams ...map[string]string) ([]byte, error)
	Del(uri string, optionalParams ...map[string]string) ([]byte, error)
}

// Client contains data needed to talk to the Trello API
type Client struct {
	host string
	key string
	token string
}

// New is a constructor for Client
func New(host string, key string, token string) Interface {
	return &Client{
		host: host,
		key: key,
		token: token,
	}
}

func (c *Client) Get(uri string, optionalParams ...map[string]string) ([]byte, error) {
	return c.request("GET", uri, optionalParams)
}

func (c *Client) Post(uri string, optionalParams ...map[string]string) ([]byte, error) {
	return c.request("POST", uri, optionalParams)
}

func (c *Client) Put(uri string, optionalParams ...map[string]string) ([]byte, error) {
	return c.request("PUT", uri, optionalParams)
}

func (c *Client) Del(uri string, optionalParams ...map[string]string) ([]byte, error) {
	return c.request("DELETE", uri, optionalParams)
}

func (c *Client) request(method string, uri string, optionalParams []map[string]string) ([]byte, error) {
	params := map[string]string{}
	if len(optionalParams) > 0 {
		params = optionalParams[0]
	}
	params["key"] = c.key
	params["token"] = c.token

	url := c.url(uri, params)
	client := http.Client{}
	req, _ := http.NewRequest(method, url, nil)

	response, err := client.Do(req)

	if err != nil {
		return nil, errors.New("Could not make " + method + " request to " + url)
	}

	body, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Response code %d", response.StatusCode)
	}

	return body, nil
}

func (c *Client) url(uri string, params map[string]string) string {
	p := url.Values{}
	for k, v := range params {
		p.Set(k, v)
	}
	query := p.Encode()

	return c.host + uri + "?" + query
}
