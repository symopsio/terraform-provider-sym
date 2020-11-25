package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
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
}

type symHttpClient struct {
	apiUrl string
}

func (c *symHttpClient) getJwt() string {
	return "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6ImtDbWVNQ0M0OGYxUDZmMTExM3hkdSJ9.eyJpc3MiOiJodHRwczovL3N5bW9wcy51cy5hdXRoMC5jb20vIiwic3ViIjoiMlIxY0lLd1Y5ekVHYldwMGd1UEVEVENCTGRpNlpnOEhAY2xpZW50cyIsImF1ZCI6Imh0dHBzOi8vYXBpLnN5bW9wcy5jb20iLCJpYXQiOjE2MDYzNDMxNDMsImV4cCI6MTYwNjQyOTU0MywiYXpwIjoiMlIxY0lLd1Y5ekVHYldwMGd1UEVEVENCTGRpNlpnOEgiLCJndHkiOiJjbGllbnQtY3JlZGVudGlhbHMifQ.SjjVBEkj_KuJcLKKyetLjHrDXtJmRWNqZzPDrtJVUF2VaozFSxti9cwrmXkgrDOv-HocNuVohdZFFljtjeCVnuMWLNU6Pu0MU9UPHQRvF5fObSeYq4A7-wMrDIyQwi6V9mc_1xpyYGSpYht_FKOmcaYnd5j9SFcgyrbUsvOJ6H0yFn_ErU5VAXOBpFC31wc3uafiaEB7CWazEFN4Bzicp6WTlLULMPNzJaNHg1vX7ccx-Wti6BPePjefOxjiASitTydtYIb0SkGX0O85no_Ipl2oiiMxRsA4t6hhlgaT-LuK2k9EwnO9GnOCzsnJZhQ7YeE4YiCXVdvGLi0hFASo3g"
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
	}
	return string(body), nil
}
