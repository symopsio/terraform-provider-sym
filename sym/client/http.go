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
	return "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6ImtDbWVNQ0M0OGYxUDZmMTExM3hkdSJ9.eyJpc3MiOiJodHRwczovL3N5bW9wcy51cy5hdXRoMC5jb20vIiwic3ViIjoiMlIxY0lLd1Y5ekVHYldwMGd1UEVEVENCTGRpNlpnOEhAY2xpZW50cyIsImF1ZCI6Imh0dHBzOi8vYXBpLnN5bW9wcy5jb20iLCJpYXQiOjE2MDYzNTE0NTMsImV4cCI6MTYwNjQzNzg1MywiYXpwIjoiMlIxY0lLd1Y5ekVHYldwMGd1UEVEVENCTGRpNlpnOEgiLCJndHkiOiJjbGllbnQtY3JlZGVudGlhbHMifQ.OFaMgBpbIFRwIjB67DqgQ4EfD3MhpDzZCcciYMk6YGN-IJmNgfX_Q44_yv1ZvmBr6y8hRQueHrCJy4vXyePzizIHf_TlXVOcHkhjjq4aMdvWBFUf392AGkiD_y42CgWbDanohGhd7OWb5T1PF-0lYz_enDDlym9bISwJ81lEQVKKnW3QiBNJKjH8UUZzk0D-nuzt8Z9beLZmZaGgrnI2PuM2k1d7IyIlHXiKBzlLhLoe7z6745_lMth4q5AflzE_wYRt7sAl2NblpNxP4BFZhMcIEf6bXa40RBDuz7zfg0zQYWQfFlXn4iwdp4O_u6KtSKTMJHXT2O6o9L0F13OfVg"
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
