package main

import (
	"github.com/scorredoira/email"
	"log"
	"net/smtp"
)

type Emailer interface {
	SendMessage(string, smtp.Auth, *email.Message) error
}

type Email int

func (e Email) SendMessage(addr string, auth smtp.Auth, m *email.Message) error {
	return email.Send(addr, auth, m)
}

func SendEmailwithAttachment(addr string, auth smtp.Auth, subject string, to string, from string, filename string) {
	var e Email = 0
	SendEmailwithAttachmentwithEmailer(addr, auth, subject, to, from, filename, e)
}

func SendEmailwithAttachmentwithEmailer(addr string, auth smtp.Auth, subject string, to string, from string, filename string, e Emailer) {
	m := email.NewMessage(subject, "")
	m.To = []string{to}
	m.From = from
	err := m.Attach(filename)
	if err != nil {
		log.Panic(err)
	}
	err = e.SendMessage(addr, auth, m)
	if err != nil {
		log.Panic(err)
	}
}

func SendEmail(addr string, auth smtp.Auth, subject string, to string, from string) {
	var e Email = 0
	SendEmailwithEmailer(addr, auth, subject, to, from, e)
}
func SendEmailwithEmailer(addr string, auth smtp.Auth, subject string, to string, from string, e Emailer) {
	m := email.NewMessage(subject, "")
	m.To = []string{to}
	m.From = from
	err := e.SendMessage(addr, auth, m)
	if err != nil {
		log.Panic(err)
	}
}
