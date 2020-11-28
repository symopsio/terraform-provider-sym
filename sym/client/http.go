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

func NewSymHttpClient() SymHttpClient {
	return &symHttpClient{
		apiUrl: "http://localhost:8000/api/v1",
	}
}

type SymHttpClient interface {
	Do(method, path string, payload interface{}) (string, error)
	Create(path string, payload interface{}, result interface{}) (string, error)
	Read(path string, result interface{}) error
}

type symHttpClient struct {
	apiUrl string
}

func (c *symHttpClient) getJwt() string {
	return "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6ImtDbWVNQ0M0OGYxUDZmMTExM3hkdSJ9.eyJpc3MiOiJodHRwczovL3N5bW9wcy51cy5hdXRoMC5jb20vIiwic3ViIjoiMlIxY0lLd1Y5ekVHYldwMGd1UEVEVENCTGRpNlpnOEhAY2xpZW50cyIsImF1ZCI6Imh0dHBzOi8vYXBpLnN5bW9wcy5jb20iLCJpYXQiOjE2MDY1ODk0ODUsImV4cCI6MTYwNjY3NTg4NSwiYXpwIjoiMlIxY0lLd1Y5ekVHYldwMGd1UEVEVENCTGRpNlpnOEgiLCJndHkiOiJjbGllbnQtY3JlZGVudGlhbHMifQ.iuvEJrnxAQGZ3qt_T6a2grM-1ebKhTByjBpRVIraFXQ21uhuEFj5tG8NjS3g0yztaUK_nqJg2mOKoqQBiplImeEjyptUpS_KhpUJsElUUN5uuu-EYHxC63Xa57CMTwNnDRN_9NlZ_2ut5eEhyKJN_L2FImZisTMdRhqrGrG6JA_mzpwIvhDNN_8Tncxo8h5gvcswJo-BW7OaM30TGZzT0EY2zVDYiq_qz6v054SWgRKaIZStSDfCuTIYG2c8WeUvXGzoSs6emwnZVuZl0Nm7Vjr96QybYZkGM4vcf295po-Wtjr6S0EeO2O3n1hwDd3d6L7ikaEoepdGgTc0asvtBQ"
}

func (c *symHttpClient) getUrl(path string) string {
	base := strings.TrimRight(c.apiUrl, "/")
	return base + "/" + strings.TrimLeft(path, "/")
}

func (c *symHttpClient) Do(method string, path string, payload interface{}) (string, error) {
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
	req.Header.Set("Authorization", "Bearer " + c.getJwt())
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	} else if resp.StatusCode >= 400 {
		return "", fmt.Errorf("error response: %v\n%s", resp, string(body))
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