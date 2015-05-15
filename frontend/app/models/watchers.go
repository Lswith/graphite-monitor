package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/revel/revel"
	"log"
	"time"
)

type PeriodicWatcher struct {
	Alarmid     string
	Notifierid  string
	Interval    string
	NotifyState State
	stopchan    chan bool
}

func (n *PeriodicWatcher) Validate(v *revel.Validation) {
	v.Required(n.Alarmid)
	v.Required(n.Notifierid)
	v.Required(n.Interval)
	v.Required(n.NotifyState)
}

func (p *PeriodicWatcher) Stop() {
	close(p.stopchan)
}

func (p *PeriodicWatcher) Run(db *bolt.DB, alarmbucket string, notifierbucket string) {
	p.stopchan = make(chan bool)
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
			log.Println("timer fired")
			a := new(Alarm)
			err := db.View(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte(alarmbucket))
				v := b.Get([]byte(p.Alarmid))
				if v == nil {
					return errors.New("alarm doesn't exist")
				}
				return a.UnMarshal(v)
			})
			if err != nil {
				break
			}
			n := new(Notifier)
			err = db.View(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte(notifierbucket))
				v := b.Get([]byte(p.Notifierid))
				if v == nil {
					return errors.New("notifier doesn't exist")
				}
				return n.UnMarshal(v)
			})
			if err != nil {
				break
			}
			state, err := a.GetState()
			if err != nil {
				break
			}
			for k, v := range state {
				log.Printf("target: %s which is in state: %s\n", k, v)
				if v == p.NotifyState {
					message := fmt.Sprintf("alarm: %s contains target: %s which is in state: %s\n", p.Alarmid, k, v)
					not, err := NewNotification(message)
					if err != nil {
						continue
					}
					err = n.Notify(not)
					if err != nil {
						continue
					}
				}
			}
			timer.Reset(d)
		case _, ok := <-p.stopchan:
			if !ok {
				log.Println("recieved a stop signal")
				break D
			}
		}
	}
}

func (p *PeriodicWatcher) Marshal() ([]byte, error) {
	return json.Marshal(p)
}

func (p *PeriodicWatcher) UnMarshal(m []byte) error {
	err := json.Unmarshal(m, p)
	if err != nil {
		return err
	}
	return nil
}
