package sym

import "log"

// Client shells out to symflow
type Client interface {
	Exec(cmd string) (string, error)
}

// NewClient creates a new symflow client
func NewClient() (Client, error) {
	return &symflow{}, nil
}

type symflow struct {
}

func (s *symflow) Exec(cmd string) (string, error) {
	log.Printf("[DEBUG] symflow: %s", cmd)
	return "FOO", nil
}
