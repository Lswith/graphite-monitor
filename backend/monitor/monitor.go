package monitor

import (
	"fmt"
	"github.com/lswith/graphite-monitor/backend/alarm"
	"github.com/lswith/graphite-monitor/backend/storage"
	"time"
)

type Monitor struct {
	periodicwatchers map[string]*PeriodicWatcher
	statefulwatchers map[string]*StatefulWatcher
}

var pbucket string = "periodicwatchers"
var sbucket string = "statefulwatchers"

func CreateMonitor() (*Monitor, error) {
	m := new(Monitor)
	m.loadPeriodicWatchers()
	m.loadStatefulWatchers()
	return m, nil
}

func (m *Monitor) CreatePeriodicWatcher(alarmid string, notifierid string, interval time.Duration, state alarm.State) (string, error) {
	p := newPeriodicWatcher(alarmid, notifierid, interval, state)
	key, err := storage.AddObject(p, pbucket)
	if err != nil {
		return "", err
	}
	m.periodicwatchers[key] = p
	go p.run()
	return key, err

}

func (m *Monitor) CreateStatefulWatcher(alarmid string, notifierid string) (string, error) {
	return "", nil
}

func (m *Monitor) DeletePeriodicWatcher(id string) error {
	if v, ok := m.periodicwatchers[id]; ok {
		v.stopchan <- true
	}
	delete(m.periodicwatchers, id)
	return storage.DeleteObject(id, pbucket)
}

func (m *Monitor) DeleteStatefulWatcher(id string) error {
	return nil
}

func (m *Monitor) GetPeriodicWatcher(id string) (*PeriodicWatcher, error) {
	if k, ok := m.periodicwatchers[id]; ok {
		return k, nil
	}
	return nil, fmt.Errorf("not found")
}

func (m *Monitor) GetStatefulWatcher(id string) (*StatefulWatcher, error) {
	return nil, nil
}

func (m *Monitor) loadPeriodicWatchers() error {
	keys, err := storage.GetKeys(pbucket)
	if err != nil {
		return err
	}
	m.periodicwatchers = make(map[string]*PeriodicWatcher)
	for _, v := range keys {
		p := new(PeriodicWatcher)
		err = storage.GetObject(v, p, pbucket)
		if err != nil {
			return err
		}
		m.periodicwatchers[v] = p
		go p.run()
	}
	return nil

}

func (m *Monitor) loadStatefulWatchers() error {
	return nil
}
