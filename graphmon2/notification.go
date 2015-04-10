package main

import (
	"net/smtp"
)

type Notifier interface {
	Notify([]Notification) error
}

type EmailNotifier struct {
	To      []string
	From    string
	Auth    smtp.Auth
	Subject string
}

func (*EmailNotifier) Notify([]Notification) error {
	return nil
}

type Notification struct {
	Message string
}
