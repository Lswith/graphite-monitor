package monitor

import (
	"encoding/json"
	"fmt"
	"github.com/lswith/graphite-monitor/backend/alarm"
	"github.com/lswith/graphite-monitor/backend/notify"
	"log"
	"time"
)

type PeriodicWatcher struct {
	Alarmid     string
	Notifierid  string
	Interval    time.Duration
	NotifyState alarm.State
	stopchan    chan bool
}

func newPeriodicWatcher(alarmid string, notifierid string, interval time.Duration, state alarm.State) *PeriodicWatcher {
	p := new(PeriodicWatcher)
	p.Alarmid = alarmid
	p.Notifierid = notifierid
	p.Interval = interval
	p.NotifyState = state
	return p
}

func (p *PeriodicWatcher) run() {
	timer := time.NewTimer(p.Interval)
	for {
		select {
		case <-timer.C:
			log.Println("timer fired")
			a, err := alarm.GetAlarm(p.Alarmid)
			if err != nil {
				break
			}
			n, err := notify.GetNotifier(p.Notifierid)
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
					not, err := notify.NewNotification(message)
					if err != nil {
						continue
					}
					err = n.Notify(not)
					if err != nil {
						continue
					}
				}
			}
			timer.Reset(p.Interval)
		case <-p.stopchan:
			break
		}
	}
}

func (p *PeriodicWatcher) Marshal() ([]byte, error) {
	return json.Marshal(p)
}

func (p *PeriodicWatcher) UnMarshal(m []byte) error {
	return json.Unmarshal(m, p)
}

type StatefulWatcher struct {
	Alarmid    string
	Notifierid string
}

func newStatefulWatcher(alarmid string, notifierid string) *StatefulWatcher {
	s := new(StatefulWatcher)
	s.Alarmid = alarmid
	s.Notifierid = notifierid
	return s
}

func (s *StatefulWatcher) run() {

}
