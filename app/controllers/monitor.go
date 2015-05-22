package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lswith/graphite-monitor/app/models"
	"github.com/revel/revel"
	"log"
	"time"
)

var (
	Periodicwatchersmap map[string]*models.PeriodicWatcher
	Statefulwatchersmap map[string]*models.StatefulWatcher
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

func (c Monitor) CreatePeriodicWatcher() revel.Result {
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
	RunPeriodicWatcher(key)
	return c.RenderText(key)
}

func (c Monitor) ReadPeriodicWatchers() revel.Result {
	log.Printf("len of Pmap: %d\n", len(Periodicwatchersmap))
	return c.RenderJson(Periodicwatchersmap)
}

func (c Monitor) ReadPeriodicWatcher(id string) revel.Result {
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

func RunPeriodicWatcher(id string) error {
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
			revel.ERROR.Println(err)
			return
		}
		timer := time.NewTimer(d)
	D:
		for {
			revel.INFO.Printf("checking state for watcher: %s\n", id)
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
						amessage := fmt.Sprintf("alarm: %s\n", p.Alarmid)
						tmessage := fmt.Sprintf("target: %s\n", k)
						smessage := fmt.Sprintf("state: %s\n", v)
						urlmessage := a.Endpoint + "/render?" + GetUrlValues(a, string(k)).Encode()
						message := amessage + tmessage + smessage + urlmessage
						revel.INFO.Println(message)
						subject := fmt.Sprintf("%s: %s", k, v)
						not, err := models.NewNotification(subject, message)
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

func (c Monitor) parseStatefulWatcher() (*models.StatefulWatcher, error) {
	statefulwatcher := new(models.StatefulWatcher)
	err := json.NewDecoder(c.Request.Body).Decode(statefulwatcher)
	return statefulwatcher, err
}

func (c Monitor) CreateStatefulWatcher() revel.Result {
	statefulwatcher, err := c.parseStatefulWatcher()
	if err != nil {
		return c.RenderError(err)
	}
	statefulwatcher.Validate(c.Validation)
	if c.Validation.HasErrors() {
		return c.RenderError(errors.New("validation error occured"))
	}
	key, err := AddObject(statefulwatcher, StatefulWatcherBucket)
	if err != nil {
		return c.RenderError(err)
	}
	Statefulwatchersmap[key] = statefulwatcher
	stopchan := make(chan bool)
	RunningWatchersmap[key] = stopchan
	RunStatefulWatcher(key)
	return c.RenderText(key)
}

func (c Monitor) ReadStatefulWatchers() revel.Result {
	log.Printf("len of Smap: %d\n", len(Statefulwatchersmap))
	return c.RenderJson(Statefulwatchersmap)
}

func (c Monitor) ReadStatefulWatcher(id string) revel.Result {
	if v, ok := Statefulwatchersmap[id]; ok {
		return c.RenderJson(v)
	}
	return c.NotFound("%s not found", id)
}

func (c Monitor) DeleteStatefulWatcher(id string) revel.Result {
	if _, ok := Statefulwatchersmap[id]; ok {
		StopWatcher(id)
		delete(Statefulwatchersmap, id)
		err := DeleteObject(id, StatefulWatcherBucket)
		if err != nil {
			return c.RenderError(err)
		}
		return c.RenderText("SUCCESS")
	}
	return c.NotFound("%s not found", id)
}

func RunStatefulWatcher(id string) error {
	if _, ok := Statefulwatchersmap[id]; !ok {
		return errors.New(fmt.Sprintf("%s not found", id))
	}
	if _, ok := RunningWatchersmap[id]; !ok {
		return errors.New(fmt.Sprintf("%s not found", id))
	}
	s := Statefulwatchersmap[id]
	stopchan := RunningWatchersmap[id]
	currstate := make(AlarmState)
	go func() {
		timer := time.NewTimer(time.Second * 10)
	D:
		for {
			select {
			case <-timer.C:
				a, err := getAlarm(s.Alarmid)
				if err != nil {
					revel.ERROR.Println(err)
					break
				}
				n, err := getNotifier(s.Notifierid)
				if err != nil {
					revel.ERROR.Println(err)
				}
				state, err := GetState(a)
				if err != nil {
					revel.ERROR.Println(err)
				}
				for k, v := range state {
					if v2, ok := currstate[k]; ok {
						if v2 == v {
							continue
						}
					}
					amessage := fmt.Sprintf("alarm: %s\n", s.Alarmid)
					tmessage := fmt.Sprintf("target: %s\n", k)
					smessage := fmt.Sprintf("state: %s\n", v)
					urlmessage := a.Endpoint + "/render?" + GetUrlValues(a, string(k)).Encode()
					message := amessage + tmessage + smessage + urlmessage
					revel.INFO.Println(message)
					subject := fmt.Sprintf("%s: %s", k, v)
					not, err := models.NewNotification(subject, message)
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
				currstate = state
				timer.Reset(time.Second * 10)
			case _, ok := <-stopchan:
				if !ok {
					break D
				}
			}
		}
	}()
	return nil
}

func getAlarm(id string) (*models.Alarm, error) {
	a := new(models.Alarm)
	err := GetObject(id, a, AlarmBucket)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func getNotifier(id string) (*models.Notifier, error) {
	n := new(models.Notifier)
	err := GetObject(id, n, NotifierBucket)
	if err != nil {
		return nil, err
	}
	return n, nil
}
