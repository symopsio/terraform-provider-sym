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
	// TODO: change this to just "symflow" when the new release is live
	exe := "/home/rory/sym/misc/cli/symflow/.venv/bin/symflow"
	svc := NewSymflowService(exe)
	_, err := svc.GetConfigValue("org")

	if err != nil {
		t.Log(err)
		t.Fail()
	}

}
