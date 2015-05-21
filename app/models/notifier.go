package models

import (
	"encoding/json"
	"github.com/revel/revel"
	"regexp"
)

type Notification struct {
	Subject string
	Body    string
}

func NewNotification(subject string, body string) (Notification, error) {
	return Notification{subject, body}, nil
}

type Notifier struct {
	From         string
	To           string
	Smtphost     string
	Smtpuser     string
	Smtppassword string
	Smtpport     string
}

var portregex = regexp.MustCompile("^[0-9]+$")

func (n *Notifier) Validate(v *revel.Validation) {
	v.Required(n.From)
	v.Email(n.From)
	v.Required(n.To)
	v.Email(n.To)
	v.Required(n.Smtphost)
	v.Match(n.Smtphost, urlregex)
	v.Required(n.Smtpuser)
	v.Required(n.Smtppassword)
	v.Required(n.Smtpport)
	v.Match(n.Smtpport, portregex)
}

func (n *Notifier) Marshal() ([]byte, error) {
	return json.Marshal(n)
}

func (n *Notifier) UnMarshal(m []byte) error {
	return json.Unmarshal(m, n)
}
