package main

import (
	"fmt"
	"github.com/lswith/graphite-monitor/alarm"
	"time"
)

func main() {
	a := new(alarm.Alarm)
	a.Name = "Test"
	fmt.Println("adding in main")
	alarmmanager.AddAlarm(a)
	time.Sleep(time.Second * 10)
	fmt.Println("deleting in main")
	alarmmanager.DeleteAlarm(a)
	time.Sleep(time.Second * 3)
}
