package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Settings map[string]string

type SymIntegration struct {
	Id             string   `json:"id"`
	OrganizationId string   `json:"organization_id"`
	Type           string   `json:"type"`
	Settings       Settings `json:"settings"`
}

type IntegrationClient interface {
	Create(integration SymIntegration) (string, error)
	Read(id string) (*SymIntegration, error)
	Update(id string, integration SymIntegration) error
	Delete(id string) error
}

func NewIntegrationClient() IntegrationClient {
	return &integrationClient{}
}

type integrationClient struct {
}

func (i *integrationClient) Create(integration SymIntegration) (string, error) {
	log.Printf("Creating integration: %v", integration)

	jwt := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6ImtDbWVNQ0M0OGYxUDZmMTExM3hkdSJ9.eyJpc3MiOiJodHRwczovL3N5bW9wcy51cy5hdXRoMC5jb20vIiwic3ViIjoiMlIxY0lLd1Y5ekVHYldwMGd1UEVEVENCTGRpNlpnOEhAY2xpZW50cyIsImF1ZCI6Imh0dHBzOi8vYXBpLnN5bW9wcy5jb20iLCJpYXQiOjE2MDYzNDMxNDMsImV4cCI6MTYwNjQyOTU0MywiYXpwIjoiMlIxY0lLd1Y5ekVHYldwMGd1UEVEVENCTGRpNlpnOEgiLCJndHkiOiJjbGllbnQtY3JlZGVudGlhbHMifQ.SjjVBEkj_KuJcLKKyetLjHrDXtJmRWNqZzPDrtJVUF2VaozFSxti9cwrmXkgrDOv-HocNuVohdZFFljtjeCVnuMWLNU6Pu0MU9UPHQRvF5fObSeYq4A7-wMrDIyQwi6V9mc_1xpyYGSpYht_FKOmcaYnd5j9SFcgyrbUsvOJ6H0yFn_ErU5VAXOBpFC31wc3uafiaEB7CWazEFN4Bzicp6WTlLULMPNzJaNHg1vX7ccx-Wti6BPePjefOxjiASitTydtYIb0SkGX0O85no_Ipl2oiiMxRsA4t6hhlgaT-LuK2k9EwnO9GnOCzsnJZhQ7YeE4YiCXVdvGLi0hFASo3g"
	integration.OrganizationId = "2R1cIKwV9zEGbWp0guPEDTCBLdi6Zg8H"
	url := "http://localhost:8000/api/v1/integrations/"
	b, err := json.Marshal(integration)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
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
	}
	log.Printf("got response: %s", string(body))

	return "", fmt.Errorf("not implemented")
}

func (i *integrationClient) Read(id string) (*SymIntegration, error) {
	panic("implement me")
}

func (i *integrationClient) Update(id string, integration SymIntegration) error {
	panic("implement me")
}

func (i *integrationClient) Delete(id string) error {
	panic("implement me")
}
