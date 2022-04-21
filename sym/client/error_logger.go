package client

import (
	"fmt"
	"log"
)

type ErrorLogger struct {
	Id            string `json:"id,omitempty"`
	IntegrationId string `json:"integration_id"`
	Destination   string `json:"destination"`
}

type ErrorLoggerClient interface {
	Create(errorLogger ErrorLogger) (string, error)
	Read(id string) (*ErrorLogger, error)
	Find(slug string) (*ErrorLogger, error)
	Update(errorLogger ErrorLogger) (string, error)
	Delete(id string) (string, error)
}

func NewErrorLoggerClient(httpClient SymHttpClient) ErrorLoggerClient {
	return &errorLoggerClient{
		HttpClient: httpClient,
	}
}

type errorLoggerClient struct {
	HttpClient SymHttpClient
}

func (c *errorLoggerClient) Create(errorLogger ErrorLogger) (string, error) {
	log.Printf("Creating ErrorLogger: %v", errorLogger)
	result := ErrorLogger{}

	if _, err := c.HttpClient.Create("/entities/error-loggers", &errorLogger, &result); err != nil {
		return "", err
	}

	if result.Id == "" {
		return "", fmt.Errorf("response indicates ErrorLogger was not created")
	}

	log.Printf("Created ErrorLogger: %s", result.Id)
	return result.Id, nil
}

func (c *errorLoggerClient) Read(id string) (*ErrorLogger, error) {
	log.Printf("Getting ErrorLogger: %s", id)
	result := ErrorLogger{}

	if err := c.HttpClient.Read(fmt.Sprintf("/entities/error-loggers/%s", id), &result); err != nil {
		return nil, err
	}

	log.Printf("Got ErrorLogger: %s", result.Id)
	return &result, nil
}

func (c *errorLoggerClient) Find(slug string) (*ErrorLogger, error) {
	log.Printf("Getting ErrorLogger by slug: %s", slug)
	var result []ErrorLogger

	if err := c.HttpClient.Read(fmt.Sprintf("/entities/error-loggers?slug=%s", slug), &result); err != nil {
		return nil, err
	}

	if len(result) != 1 {
		return nil, fmt.Errorf("one Error Logger with the slug %s was expected, but %v were found", slug, len(result))
	}

	log.Printf("Got Error Logger by slug: %s (%s)", slug, result[0].Id)
	return &result[0], nil
}

func (c *errorLoggerClient) Update(errorLogger ErrorLogger) (string, error) {
	log.Printf("Updating ErrorLogger: %v", errorLogger)
	result := ErrorLogger{}

	if _, err := c.HttpClient.Update(fmt.Sprintf("/entities/error-loggers/%s", errorLogger.Id), &errorLogger, &result); err != nil {
		return "", err
	}

	if result.Id == "" {
		return "", fmt.Errorf("response indicates ErrorLogger was not updated")
	}

	log.Printf("Updated ErrorLogger: %s", result.Id)
	return result.Id, nil
}

func (c *errorLoggerClient) Delete(id string) (string, error) {
	log.Printf("Deleting ErrorLogger: %s", id)

	if err := c.HttpClient.Delete(fmt.Sprintf("/entities/error-loggers/%s", id)); err != nil {
		return "", err
	}

	return id, nil
}
