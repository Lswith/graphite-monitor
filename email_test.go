package main

import (
	"errors"
	"github.com/scorredoira/email"
	"io"
	"net/smtp"
	"os"
	"testing"
)

var mess *email.Message = nil

func SendMessage(addr string, auth smtp.Auth, m *email.Message) error {
	mess = m
	if addr == "fail" {
		return errors.New("this is meant to return an error")
	} else {
		return nil
	}
}

func TestSendEmailwithAttachmentwithEmailer1(t *testing.T) {
	mess = email.NewMessage("", "")
	var addr = "test"
	var auth = smtp.PlainAuth("identity", "username", "password", "host")
	var to = "to"
	var from = "from"
	file, err := os.Create("test.txt")
	defer file.Close()
	defer os.Remove("test.txt")
	if err != nil {
		t.Fatal(err)
	}
	var subject = "subject"

	err = SendEmailwithAttachmentwithEmailer(addr, auth, subject, to, from, file.Name(), SendMessage)
	if err != nil {
		t.Error("shoudn't have returned an error")
	}
	if mess.From != "from" {
		t.Error("from not set")
	}
	if len(mess.To) > 1 {
		t.Error("too many persons to send email to")
	}
	if mess.To[0] != "to" {
		t.Error("to not set")
	}
	if _, ok := mess.Attachments["test.txt"]; !ok {
		t.Error("test.txt not attached")
	}
	if mess.Subject != "subject" {
		t.Error("subject not set")
	}
}

func TestSendEmailwithAttachmentwithEmailer2(t *testing.T) {
	var addr = "test"
	var auth = smtp.PlainAuth("identity", "username", "password", "host")
	var to = "to"
	var from = "from"
	if _, err := os.Stat("test.txt"); err == nil {
		os.Remove("test.txt")
	}
	var subject = "subject"
	err := SendEmailwithAttachmentwithEmailer(addr, auth, subject, to, from, "test.txt", SendMessage)
	if err == nil {
		t.Error("should have returned an error for attaching a non-existent file")
	}
}

func TestSendEmailwithAttachmentwithEmailer3(t *testing.T) {
	var addr = "fail"
	var auth = smtp.PlainAuth("identity", "username", "password", "host")
	var to = "to"
	var from = "from"
	file, err := os.Create("test.txt")
	defer file.Close()
	defer os.Remove("test.txt")
	if err != nil {
		t.Fatal(err)
	}
	var subject = "subject"
	err = SendEmailwithAttachmentwithEmailer(addr, auth, subject, to, from, file.Name(), SendMessage)
	if err == nil {
		t.Error("should have returned an error for failing to send email")
	}
}

func TestSendEmailwithEmailer1(t *testing.T) {
	mess = email.NewMessage("", "")
	var addr = "test"
	var auth = smtp.PlainAuth("identity", "username", "password", "host")
	var to = "to"
	var from = "from"
	var subject = "subject"
	err := SendEmailwithEmailer(addr, auth, subject, to, from, SendMessage)
	if err != nil {
		t.Error("shouldn't have returned an error")
	}
	if mess.From != "from" {
		t.Error("from not set")
	}
	if len(mess.To) > 1 {
		t.Error("too many persons to send email to")
	}
	if mess.To[0] != "to" {
		t.Error("to not set")
	}
	if mess.Subject != "subject" {
		t.Error("subject not set")
	}
}

func TestSendEmailwithEmailer2(t *testing.T) {
	mess = email.NewMessage("", "")
	var addr = "fail"
	var auth = smtp.PlainAuth("identity", "username", "password", "host")
	var to = "to"
	var from = "from"
	var subject = "subject"
	err := SendEmailwithEmailer(addr, auth, subject, to, from, SendMessage)
	if err == nil {
		t.Error("should have returned an error for failing to send email")
	}
}

func FakeSave1(url string, client Getter, file io.Writer) error {
	return nil
}

func FakeSave2(url string, client Getter, file io.Writer) error {
	return errors.New("save failed")
}

func FakeEmailSend1(addr string, auth smtp.Auth, subject string, to string, from string, filename string) error {
	return nil
}

func FakeEmailSend2(addr string, auth smtp.Auth, subject string, to string, from string, filename string) error {
	return errors.New("failed to send email")
}

func TestAlarmByEmail1(t *testing.T) {
	filename := "test.png"
	var a Alarm = Alarm{}
	a.Rule = "=="
	a.Target = "examples"
	a.Threshold = 0.0
	config := Config{}
	err := AlarmByEmail(a, config, filename, FakeEmailSend1, FakeSave1)
	if err != nil {
		t.Error("shouldn't have returned an error")
	}
}

func TestAlarmByEmail2(t *testing.T) {
	filename := "test.png"
	var a Alarm = Alarm{}
	a.Rule = "=="
	a.Target = "examples"
	a.Threshold = 0.0
	config := Config{}
	err := AlarmByEmail(a, config, filename, FakeEmailSend2, FakeSave1)
	if err == nil {
		t.Error("should have returned an error for not being able to send email")
	}
}

func TestAlarmByEmail3(t *testing.T) {
	filename := "test.png"
	var a Alarm = Alarm{}
	a.Rule = "=="
	a.Target = "examples"
	a.Threshold = 0.0
	config := Config{}
	err := AlarmByEmail(a, config, filename, FakeEmailSend1, FakeSave2)
	if err == nil {
		t.Error("should have returned an error for not being able to save file")
	}
}
