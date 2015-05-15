package controllers

import (
	"encoding/json"
	"errors"
	"github.com/lswith/graphite-monitor/frontend/app/models"
	"github.com/revel/revel"
)

type Notifiers struct {
	BoltController
}

func (c Notifiers) parseNotifier() (*models.Notifier, error) {
	notifier := new(models.Notifier)
	err := json.NewDecoder(c.Request.Body).Decode(notifier)
	return notifier, err
}

func (c Notifiers) Add() revel.Result {
	notifier, err := c.parseNotifier()
	if err != nil {
		return c.RenderError(err)
	}
	notifier.Validate(c.Validation)
	if c.Validation.HasErrors() {
		return c.RenderError(errors.New("validation error occured"))
	}
	b := c.Txn.Bucket([]byte(NotifierBucket))
	key, err := c.GenerateKey()
	if err != nil {
		return c.RenderError(err)
	}
	value, err := notifier.Marshal()
	if err != nil {
		return c.RenderError(err)
	}
	err = b.Put([]byte(key), value)
	if err != nil {
		return c.RenderError(err)
	}
	return c.RenderText(key)
}

func (c Notifiers) Map() revel.Result {
	m := make(map[string]*models.Notifier)
	b := c.Txn.Bucket([]byte(NotifierBucket))
	cursor := b.Cursor()
	for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
		notifier := new(models.Notifier)
		err := notifier.UnMarshal(v)
		if err != nil {
			return c.RenderError(err)
		}
		m[string(k)] = notifier
	}
	return c.RenderJson(m)
}

func (c Notifiers) Get(id string) revel.Result {
	notifier := new(models.Notifier)
	b := c.Txn.Bucket([]byte(NotifierBucket))
	v := b.Get([]byte(id))
	if v == nil {
		return c.NotFound("%s not found", id)
	}
	err := notifier.UnMarshal(v)
	if err != nil {
		return c.RenderError(err)
	}
	return c.RenderJson(notifier)
}

func (c Notifiers) Delete(id string) revel.Result {
	b := c.Txn.Bucket([]byte(NotifierBucket))
	err := b.Delete([]byte(id))
	if err != nil {
		return c.RenderError(err)
	}
	return c.RenderText("SUCCESS")
}
