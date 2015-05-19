package models

import (
	"encoding/json"
	"github.com/revel/revel"
)

type Notification string

func NewNotification(message string) (Notification, error) {
	return Notification(message), nil
}

type Notifier struct {
}

func (n *Notifier) Validate(v *revel.Validation) {
}

func (n *Notifier) Marshal() ([]byte, error) {
	return json.Marshal(n)
}

func (n *Notifier) UnMarshal(m []byte) error {
	return json.Unmarshal(m, n)
}
