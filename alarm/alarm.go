package alarm

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Alarm struct {
	Name string
}

func (a *Alarm) Run() State {
	if rand.Intn(2) == 0 {
		return false
	}
	return true
}

func AddAlarm(a *Alarm) {
	fmt.Println("trying to add alarm")
	add <- a
}

func DeleteAlarm(a *Alarm) {
	fmt.Println("trying to delete alarm")
	del <- a
}

func GetAlarms() map[*Alarm]State {
	alarmslock.Lock()
	m := make(map[*Alarm]State)
	for k, v := range alarms {
		m[k] = v
	}
	alarmslock.Unlock()
	return m
}

type State bool

var (
	add        chan *Alarm
	del        chan *Alarm
	alarms     map[*Alarm]State
	updates    chan Update
	alarmslock *sync.Mutex
)

type Update struct {
	Current State
	A       *Alarm
}

func run() {
	for {
		select {
		case a := <-add:
			fmt.Printf("adding %s\n", a.Name)
			alarms[a] = false
		case d := <-del:
			fmt.Printf("deleting %s\n", d.Name)
			delete(alarms, d)

		default:
			alarmslock.Lock()
			for alarm, oldstate := range alarms {
				var newstate = alarm.Run()
				if newstate != oldstate {
					u := Update{
						Current: newstate,
						A:       alarm,
					}
					updates <- u
					alarms[alarm] = newstate
				}
			}
			alarmslock.Unlock()
			time.Sleep(time.Second)
		}
	}
}

func init() {
	add = make(chan *Alarm)
	del = make(chan *Alarm)
	alarms = make(map[*Alarm]State)
	updates = make(chan Update)
	alarmslock = new(sync.Mutex)
	go run()
	go func() {
		for {
			select {
			case u := <-updates:
				notify(u)
			}
		}
	}()
}
