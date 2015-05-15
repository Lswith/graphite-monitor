package controllers

import (
	"encoding/json"
	"errors"
	"github.com/lswith/graphite-monitor/frontend/app/models"
	"github.com/revel/revel"
	"log"
)

var (
	Periodicwatchersmap map[string]*models.PeriodicWatcher
)

type Monitor struct {
	BoltController
}

func (c Monitor) parsePeriodicWatcher() (*models.PeriodicWatcher, error) {
	periodicwatcher := new(models.PeriodicWatcher)
	err := json.NewDecoder(c.Request.Body).Decode(periodicwatcher)
	return periodicwatcher, err
}

func (c Monitor) AddPeriodicWatcher() revel.Result {
	periodicwatcher, err := c.parsePeriodicWatcher()
	if err != nil {
		return c.RenderError(err)
	}
	periodicwatcher.Validate(c.Validation)
	if c.Validation.HasErrors() {
		return c.RenderError(errors.New("validation error occured"))
	}
	b := c.Txn.Bucket([]byte(PeriodicWatcherBucket))
	key, err := c.GenerateKey()
	if err != nil {
		return c.RenderError(err)
	}
	value, err := periodicwatcher.Marshal()
	if err != nil {
		return c.RenderError(err)
	}
	err = b.Put([]byte(key), value)
	if err != nil {
		return c.RenderError(err)
	}
	Periodicwatchersmap[key] = periodicwatcher
	go periodicwatcher.Run(Db, AlarmBucket, NotifierBucket)
	return c.RenderText(key)
}

func (c Monitor) MapPeriodicWatchers() revel.Result {
	log.Printf("len of Pmap: %d\n", len(Periodicwatchersmap))
	return c.RenderJson(Periodicwatchersmap)
}

func (c Monitor) GetPeriodicWatcher(id string) revel.Result {
	if v, ok := Periodicwatchersmap[id]; ok {
		return c.RenderJson(v)
	}
	return c.NotFound("%s not found", id)
}

func (c Monitor) DeletePeriodicWatcher(id string) revel.Result {
	if v, ok := Periodicwatchersmap[id]; ok {
		v.Stop()
		delete(Periodicwatchersmap, id)
		b := c.Txn.Bucket([]byte(PeriodicWatcherBucket))
		err := b.Delete([]byte(id))
		if err != nil {
			return c.RenderError(err)
		}
		return c.RenderText("SUCCESS")
	}
	return c.NotFound("%s not found", id)
}
