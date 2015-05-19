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
