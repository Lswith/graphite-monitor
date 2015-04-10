package main

import (
	// "fmt"
	"github.com/op/go-logging"
	"os"
	"testing"
)

func TestSetupLogging(t *testing.T) {
	if log == nil {
		t.Error("logging didn't initialize")
	}
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	if err := SetupLogging(backend); err != nil {
		t.Error("SetupLogging should not have thrown an error", err)
	}
}

func TestCheckSettings(t *testing.T) {
	if settings == nil {
		t.Error("Settings should not be nil")
	}
}

type FakeDataGetter struct {
	data Data
}

func (f *FakeDataGetter) GetDataForTarget(target string) (Data, error) {
	return f.data, nil

}

func TestGenerateNotifications1(t *testing.T) {
	alarms := make([]Alarm, 0)
	fakegetter := FakeDataGetter{}
	notifications, err := GenerateNotifications(alarms, &fakegetter)
	if err != nil {
		t.Error("GenerateNotifications shouldn't throw an error")
	}
	if len(notifications) > 0 {
		t.Error("GenerateNotifications shouldn't have returned any notifications")
	}
}

func TestGenerateNotifications2(t *testing.T) {
	alarms := make([]Alarm, 1)
	alarms[0].Threshold = 0
	alarms[0].Rule = "!="
	alarms[0].Target = "test"
	alarms[0].enabled = true
	fakegetter := FakeDataGetter{}
	fakegetter.data = Data{
		[][2]float64{
			[2]float64{
				1.1,
				2.2,
			},
		},
	}
	notifications, err := GenerateNotifications(alarms, &fakegetter)
	if err != nil {
		t.Error("GenerateNotifications shouldn't throw an error")
	}
	if len(notifications) != 1 {
		t.Error("GenerateNotifications should have returned 1 notification")
	}
	var testmsg = "Rule: " + alarms[0].Rule + " has been met for target: " + alarms[0].Target
	if notifications[0].Message != testmsg {
		t.Error("GenerateNotifications should have produced the test message: " + testmsg)
	}
}