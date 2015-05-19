package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/lswith/graphite-monitor/app/models"
	"github.com/revel/revel"
	"log"
	"time"
)

var (
	Periodicwatchersmap map[string]*models.PeriodicWatcher
	RunningWatchersmap  map[string]chan bool
)

type Monitor struct {
	*revel.Controller
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
	key, err := AddObject(periodicwatcher, PeriodicWatcherBucket)
	if err != nil {
		return c.RenderError(err)
	}
	Periodicwatchersmap[key] = periodicwatcher
	stopchan := make(chan bool)
	RunningWatchersmap[key] = stopchan
	RunWatcher(key)
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
	if _, ok := Periodicwatchersmap[id]; ok {
		StopWatcher(id)
		delete(Periodicwatchersmap, id)
		err := DeleteObject(id, PeriodicWatcherBucket)
		if err != nil {
			return c.RenderError(err)
		}
		return c.RenderText("SUCCESS")
	}
	return c.NotFound("%s not found", id)
}

func RunWatcher(id string) error {
	if _, ok := Periodicwatchersmap[id]; !ok {
		return errors.New(fmt.Sprintf("%s not found", id))
	}
	if _, ok := RunningWatchersmap[id]; !ok {
		return errors.New(fmt.Sprintf("%s not found", id))
	}
	p := Periodicwatchersmap[id]
	stopchan := RunningWatchersmap[id]
	go func() {
		d, err := time.ParseDuration(p.Interval)
		if err != nil {
			log.Println(err)
			return
		}
		timer := time.NewTimer(d)
	D:
		for {
			select {
			case <-timer.C:
				a, err := getAlarm(p.Alarmid)
				if err != nil {
					revel.ERROR.Println(err)
					break
				}
				n, err := getNotifier(p.Notifierid)
				if err != nil {
					revel.ERROR.Println(err)
				}
				state, err := GetState(a)
				if err != nil {
					revel.ERROR.Println(err)
				}
				for k, v := range state {
					if v == p.NotifyState {
						message := fmt.Sprintf("alarm: %s contains target: %s which is in state: %s\n", p.Alarmid, k, v)
						not, err := models.NewNotification(message)
						if err != nil {
							revel.ERROR.Println(err)
							continue
						}
						err = Notify(n, not)
						if err != nil {
							revel.ERROR.Println(err)
							continue
						}
					}
				}
				timer.Reset(d)
			case _, ok := <-stopchan:
				if !ok {
					break D
				}
			}
		}
	}()
	return nil
}

func StopWatcher(id string) error {
	if v, ok := RunningWatchersmap[id]; ok {
		close(v)
		return nil
	}
	return errors.New(fmt.Sprintf("%s not found", id))
}

func getAlarm(alarmid string) (*models.Alarm, error) {
	a := new(models.Alarm)
	err := GetObject(id, a, AlarmBucket)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func getNotifier(notifierid string) (*models.Notifier, error) {
	n := new(models.Notifier)
	err := GetObject(id, notifier, NotifierBucket)
	if err != nil {
		return nil, err
	}
	return n, nil
}
