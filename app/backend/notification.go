package backend

import (
	"fmt"
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

type ConsoleNotifier struct {
}

func (*ConsoleNotifier) Notify(notifications []Notification) error {
	for _, notification := range notifications {
		fmt.Println(notification.Message)
	}
	return nil
}
