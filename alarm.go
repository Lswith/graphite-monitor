package main

import (
	// "fmt"
	"math/rand"
	"time"
)

type alarm struct {
	snooze   chan time.Duration
	state    chan bool
	alm      Alm
	timer    *time.Timer
	curState bool
}

type Alm struct {
	stuff string
}

func (a *alarm) run() {
L:
	for {

		select {
		case d, ok := <-a.snooze:
			if ok {
				a.timer.Reset(d)
			} else {
				close(a.state)
				break L
			}
		case <-a.timer.C:
			s1 := a.GetNewState()
			if s1 != a.curState {
				a.curState = s1
				a.state <- s1
			}
			a.timer.Reset(time.Second)
		}
	}
}

func (a *alarm) GetNewState() bool {
	if rand.Intn(2) == 0 {
		return true
	}
	return false
}

type alarmManager struct {
	alarms []*alarm
}

func (am *alarmManager) AddNewAlarm(a Alm) <-chan bool {
	alarm := new(alarm)
	alarm.alm = a
	alarm.snooze = make(chan time.Duration, 20)
	alarm.state = make(chan bool, 256)
	alarm.timer = time.NewTimer(time.Second)
	am.alarms = append(am.alarms, alarm)
	go alarm.run()
	return alarm.state
}

func (am *alarmManager) DeleteAlarm(a Alm) {
	for i, v := range am.alarms {
		if v.alm == a {
			close(v.snooze)
			am.alarms[i] = am.alarms[len(am.alarms)-1]
			am.alarms[len(am.alarms)-1] = nil
			am.alarms = am.alarms[:len(am.alarms)-1]
			break
		}
	}

}

func (am *alarmManager) SnoozeAlarm(a Alm, d time.Duration) {
	for _, v := range am.alarms {
		if v.alm == a {
			v.snooze <- d
			break
		}
	}
}
