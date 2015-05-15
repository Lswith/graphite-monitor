package controllers

import (
	"encoding/json"
	"errors"
	"github.com/lswith/graphite-monitor/frontend/app/models"
	"github.com/revel/revel"
)

type Alarms struct {
	BoltController
}

func (c Alarms) parseAlarm() (*models.Alarm, error) {
	alarm := new(models.Alarm)
	err := json.NewDecoder(c.Request.Body).Decode(alarm)
	return alarm, err
}

func (c Alarms) Add() revel.Result {
	alarm, err := c.parseAlarm()
	if err != nil {
		return c.RenderError(err)
	}
	alarm.Validate(c.Validation)
	if c.Validation.HasErrors() {
		return c.RenderError(errors.New("validation error occured"))
	}
	b := c.Txn.Bucket([]byte(AlarmBucket))
	key, err := c.GenerateKey()
	if err != nil {
		return c.RenderError(err)
	}
	value, err := alarm.Marshal()
	if err != nil {
		return c.RenderError(err)
	}
	err = b.Put([]byte(key), value)
	if err != nil {
		return c.RenderError(err)
	}
	return c.RenderText(key)
}

func (c Alarms) Map() revel.Result {
	m := make(map[string]*models.Alarm)
	b := c.Txn.Bucket([]byte(AlarmBucket))
	cursor := b.Cursor()
	for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
		alarm := new(models.Alarm)
		err := alarm.UnMarshal(v)
		if err != nil {
			return c.RenderError(err)
		}
		m[string(k)] = alarm
	}
	return c.RenderJson(m)
}

func (c Alarms) Get(id string) revel.Result {
	alarm := new(models.Alarm)
	b := c.Txn.Bucket([]byte(AlarmBucket))
	v := b.Get([]byte(id))
	if v == nil {
		return c.NotFound("%s not found", id)
	}
	err := alarm.UnMarshal(v)
	if err != nil {
		return c.RenderError(err)
	}
	return c.RenderJson(alarm)
}

func (c Alarms) Delete(id string) revel.Result {
	b := c.Txn.Bucket([]byte(AlarmBucket))
	err := b.Delete([]byte(id))
	if err != nil {
		return c.RenderError(err)
	}
	return c.RenderText("SUCCESS")
}
