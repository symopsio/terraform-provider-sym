package service

import (
	"testing"
)

func TestGetVersion(t *testing.T) {
	svc := NewSymflowService()
	_, err := svc.GetVersion()

	if err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestGetValidConfigValue(t *testing.T) {
	svc := NewSymflowService()
	_, err := svc.GetConfigValue("org")

	if err != nil {
		t.Log(err)
		t.Fail()
	}

}
