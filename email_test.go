package main

import (
	"errors"
	"github.com/scorredoira/email"
	"net/smtp"
	"os"
	"testing"
)

type TestEmail int

var mess *email.Message

func (e TestEmail) SendMessage(addr string, auth smtp.Auth, m *email.Message) error {
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
	var e TestEmail = 0
	SendEmailwithAttachmentwithEmailer(addr, auth, subject, to, from, file.Name(), e)
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
	recovered := false
	defer func() {
		if r := recover(); r != nil {
			recovered = true
		}
	}()
	var addr = "test"
	var auth = smtp.PlainAuth("identity", "username", "password", "host")
	var to = "to"
	var from = "from"
	if _, err := os.Stat("test.txt"); err == nil {
		os.Remove("test.txt")
	}
	var subject = "subject"
	var e TestEmail = 0
	SendEmailwithAttachmentwithEmailer(addr, auth, subject, to, from, "test.txt", e)
	if recovered == false {
		t.Error("should have paniced trying to attach a non-existent file")
	}
}

func TestSendEmailwithAttachmentwithEmailer3(t *testing.T) {
	recovered := false
	defer func() {
		if r := recover(); r != nil {
			recovered = true
		}
	}()
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
	var e TestEmail = 0
	SendEmailwithAttachmentwithEmailer(addr, auth, subject, to, from, file.Name(), e)
	if recovered == false {
		t.Error("should have paniced trying to send email")
	}
}

func TestSendEmailwithEmailer1(t *testing.T) {
	mess = email.NewMessage("", "")
	var addr = "test"
	var auth = smtp.PlainAuth("identity", "username", "password", "host")
	var to = "to"
	var from = "from"
	var subject = "subject"
	var e TestEmail = 0
	SendEmailwithEmailer(addr, auth, subject, to, from, e)
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
	recovered := false
	defer func() {
		if r := recover(); r != nil {
			recovered = true
		}
	}()
	var addr = "fail"
	var auth = smtp.PlainAuth("identity", "username", "password", "host")
	var to = "to"
	var from = "from"
	var subject = "subject"
	var e TestEmail = 0
	SendEmailwithEmailer(addr, auth, subject, to, from, e)
	if recovered == false {
		t.Error("should have paniced trying to send email")
	}
}
