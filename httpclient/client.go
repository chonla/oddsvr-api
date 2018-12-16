package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Client struct {
	AccessToken string
}

func NewClient() *Client {
	return &Client{}
}

func NewClientWithToken(token string) *Client {
	return &Client{
		AccessToken: token,
	}
}

func (c *Client) Get(url string, output interface{}) error {
	httpClient := &http.Client{}
	req, e := http.NewRequest("GET", url, nil)
	if e != nil {
		return e
	}
	if c.AccessToken != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))
	}

	resp, e := httpClient.Do(req)
	if e != nil {
		return e
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		e = json.Unmarshal(bodyBytes, output)
		if e != nil {
			return e
		}
		return nil
	}
	return fmt.Errorf("error: %s", resp.Status)
}

func (c *Client) Post(url string, data, output interface{}) error {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(data)

	httpClient := &http.Client{}
	req, e := http.NewRequest("POST", url, b)
	if e != nil {
		return e
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", fmt.Sprintf("%d", len(b.String())))
	if c.AccessToken != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))
	}

	resp, e := httpClient.Do(req)
	if e != nil {
		return e
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		e = json.Unmarshal(bodyBytes, output)
		if e != nil {
			return e
		}
		return nil
	}
	return fmt.Errorf("error: %s", resp.Status)
}

func (c *Client) PostForm(url, data string, output interface{}) error {
	httpClient := &http.Client{}
	req, e := http.NewRequest("POST", url, strings.NewReader(data))
	if e != nil {
		return e
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", fmt.Sprintf("%d", len(data)))
	if c.AccessToken != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))
	}

	resp, e := httpClient.Do(req)
	if e != nil {
		return e
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		e = json.Unmarshal(bodyBytes, output)
		if e != nil {
			return e
		}
		return nil
	}
	return fmt.Errorf("error: %s", resp.Status)
}
