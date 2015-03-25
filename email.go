package main

import (
	"github.com/scorredoira/email"
	"log"
	"net/smtp"
)

func sendEmailwithAttachment(addr string, auth smtp.Auth, subject string, to string, from string, filename string) {
	m := email.NewMessage(subject, "")
	m.To = []string{to}
	m.From = from
	err := m.Attach(filename)
	if err != nil {
		log.Panic(err)
	}
	err = email.Send(addr, auth, m)
	if err != nil {
		log.Panic(err)
	}
}

func sendEmail(addr string, auth smtp.Auth, subject string, to string, from string) {
	m := email.NewMessage(subject, "")
	m.To = []string{to}
	m.From = from
	err := email.Send(addr, auth, m)
	if err != nil {
		log.Panic(err)
	}
}
