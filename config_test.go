package main

import (
	"strings"
	"testing"
)

const example1 = `
{
	"endpoint":"http://example.com",
	"interval":"-20mins",
	"target":"test",
	"threshold":0.0,
	"frequency":"20m",
	"rule":"<=",
	"emailserver":"smtp.example.com",
	"emailto":"to@example.com",
	"emailfrom":"from@example.com",
	"emailuser":"user",
	"emailpassword":"password",
    "emailport":"587",
    "emailsubject":"Alert"
}
`

func TestReadConfig1(t *testing.T) {

	c1 := ReadConfig(strings.NewReader(example1))
	if c1.Endpoint != "http://example.com" {
		t.Error("endpoint not being read correctly")
	}
	if c1.Interval != "-20mins" {
		t.Error("interval not being read correctly")
	}
	if c1.Target != "test" {
		t.Error("target not being read correctly")
	}
	if c1.Threshold != 0.0 {
		t.Error("threshold not being read correctly")
	}
	if c1.Frequency != "20m" {
		t.Error("frequency not being read correctly")
	}
	if c1.Rule != "<=" {
		t.Error("rule not being read correctly")
	}
	if c1.EmailServer != "smtp.example.com" {
		t.Error("emailserver not being read correctly")
	}
	if c1.EmailSubject != "Alert" {
		t.Error("email subject not being read correctly")
	}
	if c1.EmailPort != "587" {
		t.Error("email port not being read correctly")
	}
	if c1.EmailPassword != "password" {
		t.Error("email password not being read correctly")
	}
	if c1.EmailFrom != "from@example.com" {
		t.Error("email from not being read correctly")
	}
	if c1.EmailUser != "user" {
		t.Error("email user not being read correctly")
	}
	if c1.EmailTo != "to@example.com" {
		t.Error("email to not being read correctly")
	}
}

func TestReadConfig2(t *testing.T) {
	recovered := false
	defer func() {
		if r := recover(); r != nil {
			recovered = true
		}
	}()
	ReadConfig(strings.NewReader("hello"))
	t.Error("should have paniced")
}
