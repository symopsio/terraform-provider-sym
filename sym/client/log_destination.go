package client

import (
	"fmt"
	"log"
)

type LogDestination struct {
	Id            string   `json:"id,omitempty"`
	Type          string   `json:"type"`
	IntegrationId string   `json:"integration_id"`
	Settings      Settings `json:"settings"`
}

func (s LogDestination) String() string {
	return fmt.Sprintf("{id=%s, type=%s, settings=%v}", s.Id, s.Type, s.Settings)
}

type LogDestinationClient interface {
	Create(destination LogDestination) (string, error)
	Read(id string) (*LogDestination, error)
	Find(name, destinationType string) (*LogDestination, error)
	Update(destination LogDestination) (string, error)
	Delete(id string) (string, error)
}

func NewLogDestinationClient(httpClient SymHttpClient) LogDestinationClient {
	return &logDestinationClient{
		HttpClient: httpClient,
	}
}

type logDestinationClient struct {
	HttpClient SymHttpClient
}

func (l *logDestinationClient) Create(destination LogDestination) (string, error) {
	log.Printf("Creating Sym LogDestination: %v", destination)

	result := LogDestination{}
	if _, err := l.HttpClient.Create("/entities/log-destinations", &destination, &result); err != nil {
		return "", err
	}

	if result.Id == "" {
		return "", fmt.Errorf("response indicates Sym LogDestination was not created")
	}

	log.Printf("Created Sym LogDestination: %s", result.Id)
	return result.Id, nil
}

func (l *logDestinationClient) Read(id string) (*LogDestination, error) {
	log.Printf("Getting Sym LogDestination: %s", id)
	result := LogDestination{}

	if err := l.HttpClient.Read(fmt.Sprintf("/entities/log-destinations/%s", id), &result); err != nil {
		return nil, err
	}

	log.Printf("Got Sym LogDestination: %s", result.Id)
	return &result, nil
}

func (l *logDestinationClient) Find(name, destinationType string) (*LogDestination, error) {
	log.Printf("Getting Sym Log Destination by type and name: %s %s", destinationType, name)
	var result []LogDestination

	if err := l.HttpClient.Read(fmt.Sprintf("/entities/log-destinations?slug=%s&type=%s", name, destinationType), &result); err != nil {
		return nil, err
	}

	if len(result) != 1 {
		return nil, fmt.Errorf("one Log Destination of type %s with the name %s was expected, but %v were found", destinationType, name, len(result))
	}

	log.Printf("Got Sym Log Destination by type and name: %s %s (%s)", destinationType, name, result[0].Id)
	return &result[0], nil
}

func (l *logDestinationClient) Update(destination LogDestination) (string, error) {
	log.Printf("Updating Sym LogDestination: %v", destination)
	result := LogDestination{}

	if _, err := l.HttpClient.Update(fmt.Sprintf("/entities/log-destinations/%s", destination.Id), &destination, &result); err != nil {
		return "", err
	}

	if result.Id == "" {
		return "", fmt.Errorf("response indicates Sym LogDestination was not updated")
	}

	log.Printf("Updated Sym LogDestination: %s", result.Id)
	return result.Id, nil
}

func (l *logDestinationClient) Delete(id string) (string, error) {
	log.Printf("Deleting Sym LogDestination: %s", id)

	if err := l.HttpClient.Delete(fmt.Sprintf("/entities/log-destinations/%s", id)); err != nil {
		return "", err
	}

	return id, nil
}
