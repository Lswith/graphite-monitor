package main

import (
	"errors"
	"os"
	"testing"
)

func TestSetup(t *testing.T) {
	file, err := os.Create("test.json")
	if err != nil {
		t.Error(err)
	}
	defer file.Close()
	defer os.Remove(file.Name())
	file.WriteString(example1)
	_, err = Setup("test.log", "test.json")
	if err != nil {
		t.Error(err)
	}
	err = os.Remove("test.log")
	if err != nil {
		t.Error(err)
	}
}

func FakeGetData1(config Config, client Getter) ([]Data, error) {
	return []Data{}, nil
}

func FakeGetData2(config Config, client Getter) ([]Data, error) {
	return []Data{}, errors.New("couldn't get data")
}

func FakeMonitorData1(d []Data, rule string, thres float64) ([]Alarm, error) {
	return []Alarm{
		Alarm{
			"example.stats",
			">",
			0.0,
		},
	}, nil
}
func FakeMonitorData2(d []Data, rule string, thres float64) ([]Alarm, error) {
	return []Alarm{}, errors.New("failed to monitor data")
}

func FakeAlarmByEmail1(alarm Alarm, config Config, filename string, emailsend SendEmailwithAttachmentFunc, save SaveFileFunc) error {
	return nil
}

func FakeAlarmByEmail2(alarm Alarm, config Config, filename string, emailsend SendEmailwithAttachmentFunc, save SaveFileFunc) error {
	return errors.New("failed to alarm by email")
}

func TestLoop1(t *testing.T) {
	config := Config{}
	err := Loop(config, FakeGetData1, FakeMonitorData1, FakeAlarmByEmail1)
	if err != nil {
		t.Error("shouldn't have returned an error")
	}
}

func TestLoop2(t *testing.T) {
	config := Config{}
	err := Loop(config, FakeGetData2, FakeMonitorData1, FakeAlarmByEmail1)
	if err == nil {
		t.Error("should have returned an error for not being able to get data")
	}
}

func TestLoop3(t *testing.T) {
	config := Config{}
	err := Loop(config, FakeGetData1, FakeMonitorData2, FakeAlarmByEmail1)
	if err == nil {
		t.Error("should have returned an error for not being able to monitor data")
	}
}

func TestLoop4(t *testing.T) {
	config := Config{}
	err := Loop(config, FakeGetData1, FakeMonitorData1, FakeAlarmByEmail2)
	if err != nil {
		t.Error("shouldn't have returned an error for not being able to alarm by email")
	}
}
