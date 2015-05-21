package controllers

import (
	"encoding/json"
	"errors"
	"github.com/lswith/graphite-monitor/app/models"
	"github.com/revel/revel"
	"github.com/scorredoira/email"
	"net/smtp"
)

type Notifiers struct {
	*revel.Controller
}

func (c Notifiers) parseNotifier() (*models.Notifier, error) {
	notifier := new(models.Notifier)
	err := json.NewDecoder(c.Request.Body).Decode(notifier)
	return notifier, err
}

func (c Notifiers) CreateNotifier() revel.Result {
	notifier, err := c.parseNotifier()
	if err != nil {
		return c.RenderError(err)
	}
	notifier.Validate(c.Validation)
	if c.Validation.HasErrors() {
		return c.RenderError(errors.New("validation error occured"))
	}
	key, err := AddObject(notifier, NotifierBucket)
	if err != nil {
		return c.RenderError(err)
	}
	return c.RenderText(key)
}

func (c Notifiers) ReadNotifiers() revel.Result {
	m := make(map[string]*models.Notifier)
	ids, err := GetKeys(NotifierBucket)
	if err != nil {
		return c.RenderError(err)
	}
	for _, id := range ids {
		notifier := new(models.Notifier)
		err = GetObject(id, notifier, NotifierBucket)
		if err != nil {
			return c.RenderError(err)
		}
		m[id] = notifier
	}
	return c.RenderJson(m)
}

func (c Notifiers) ReadNotifier(id string) revel.Result {
	notifier := new(models.Notifier)
	err := GetObject(id, notifier, NotifierBucket)
	if err != nil {
		c.RenderError(err)
	}
	return c.RenderJson(notifier)
}

func (c Notifiers) DeleteNotifier(id string) revel.Result {
	err := DeleteObject(id, NotifierBucket)
	if err != nil {
		return c.RenderError(err)
	}
	return c.RenderText("SUCCESS")
}

func Notify(n *models.Notifier, notification models.Notification) error {
	revel.INFO.Printf("sent notification: %s %s\n", notification.Subject, notification.Body)
	m := email.NewMessage(notification.Subject, notification.Body)
	m.To = []string{n.To}
	m.From = n.From
	auth := smtp.PlainAuth("", n.Smtpuser, n.Smtppassword, n.Smtphost)
	return email.Send(n.Smtphost+":"+n.Smtpport, auth, m)
}
