package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func NewSymHttpClient(apiUrl string) SymHttpClient {
	return &symHttpClient{
		apiUrl:       apiUrl,
		configReader: NewConfigReader(),
	}
}

type SymHttpClient interface {
	Do(method, path string, payload interface{}) (string, error)
	Create(path string, payload interface{}, result interface{}) (string, error)
	Read(path string, result interface{}) error
	Update(path string, payload interface{}, result interface{}) (string, error)
	Delete(path string) error
}

type symHttpClient struct {
	apiUrl       string
	configReader ConfigReader
}

func (c *symHttpClient) getJwt() (string, error) {
	return c.configReader.GetJwt()
}

func (c *symHttpClient) getUrl(path string) string {
	base := strings.TrimRight(c.apiUrl, "/")
	return base + "/" + strings.TrimLeft(path, "/")
}

func (c *symHttpClient) Do(method string, path string, payload interface{}) (string, error) {
	jwt, err := c.getJwt()
	if err != nil {
		return "", err
	}
	url := c.getUrl(path)
	b, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	log.Printf("submitting request: %s %s %s", method, path, string(b))
	req, err := http.NewRequest(method, url, bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+jwt)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	} else if resp.StatusCode >= 400 {
		return "", fmt.Errorf("Error %d: %s", resp.StatusCode, string(body))
	}
	return string(body), nil
}

func (c *symHttpClient) Create(path string, payload interface{}, result interface{}) (string, error) {
	body, err := c.Do("POST", path, payload)
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal([]byte(body), result); err != nil {
		return "", err
	}

	log.Printf("got response: %v", result)
	return body, nil
}

func (c *symHttpClient) Read(path string, result interface{}) error {
	body, err := c.Do("GET", path, nil)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(body), result)
}

func (c *symHttpClient) Update(path string, payload interface{}, result interface{}) (string, error) {
	body, err := c.Do("PATCH", path, payload)
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal([]byte(body), result); err != nil {
		return "", err
	}

	log.Printf("got response: %v", result)
	return body, nil
}

func (c *symHttpClient) Delete(path string) error {
	if _, err := c.Do("DELETE", path, nil); err != nil {
		return err
	}

	return nil
}
