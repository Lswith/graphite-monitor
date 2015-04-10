package main

import (
	"time"
)

type Settings struct {
	Frequency time.Duration
	Alarms    []Alarm
	Notifiers []Notifier
	Graphite  *GraphiteGetter
}

func (s *Settings) UpdateSettings() {

}
