package service

import (
	"testing"
)

func TestGetVersion(t *testing.T) {
	svc := NewSymflowService("symflow")
	_, err := svc.GetVersion()

	if err != nil {
		t.Fail()
	}
}

func TestGetValidConfigValue(t *testing.T) {
	exe := "symflow"
	svc := NewSymflowService(exe)
	_, err := svc.GetConfigValue("org")

	if err != nil {
		t.Log(err)
		t.Fail()
	}

}
