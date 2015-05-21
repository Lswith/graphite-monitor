package models

import (
	"encoding/json"
	"github.com/revel/revel"
)

type PeriodicWatcher struct {
	Alarmid     string
	Notifierid  string
	Interval    string
	NotifyState State
	stopchan    chan bool
}

type StatefulWatcher struct {
	Alarmid    string
	Notifierid string
	stopchan   chan bool
}

func (n *PeriodicWatcher) Validate(v *revel.Validation) {
	v.Required(n.Alarmid)
	v.Required(n.Notifierid)
	v.Required(n.Interval)
	v.Required(n.NotifyState)
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

func (n *StatefulWatcher) Validate(v *revel.Validation) {
	v.Required(n.Alarmid)
	v.Required(n.Notifierid)
}

func (p *StatefulWatcher) Marshal() ([]byte, error) {
	return json.Marshal(p)
}

func (p *StatefulWatcher) UnMarshal(m []byte) error {
	err := json.Unmarshal(m, p)
	if err != nil {
		return err
	}
	return nil
}
