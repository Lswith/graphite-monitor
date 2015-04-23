package controllers

import (
	"github.com/lswith/graphite-monitor/alarm"
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

var alarms map[string]*alarm.Alarm

func init() {
	alarms = make(map[string]*alarm.Alarm)
}

func (c App) AddAlarm(name string) revel.Result {
	a := new(alarm.Alarm)
	a.Name = name
	alarms[name] = a
	alarm.AddAlarm(a)
	return c.Render(name)
}

func (c App) DeleteAlarm(name string) revel.Result {
	a := alarms[name]
	delete(alarms, name)
	alarm.DeleteAlarm(a)
	return c.Render(name)
}
