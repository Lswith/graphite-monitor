package main

import (
	"errors"
	// "time"
)

type Alarm struct {
	Target    string
	Threshold float64
	Rule      string
	enabled   bool
}

func (a *Alarm) HasRuleBeenMet(data Data) (bool, error) {
	switch a.Rule {
	case "==":
		for j := range data.DataPoints {
			if data.DataPoints[j][0] == a.Threshold {
				return true, nil
			}
		}
	case "!=":
		for j := range data.DataPoints {
			if data.DataPoints[j][0] != a.Threshold {
				return true, nil
			}
		}
	case "<":
		for j := range data.DataPoints {
			if data.DataPoints[j][0] < a.Threshold {
				return true, nil
			}
		}
	case "<=":
		for j := range data.DataPoints {
			if data.DataPoints[j][0] <= a.Threshold {
				return true, nil
			}
		}
	case ">":
		for j := range data.DataPoints {
			if data.DataPoints[j][0] > a.Threshold {
				return true, nil
			}
		}
	case ">=":
		for j := range data.DataPoints {
			if data.DataPoints[j][0] >= a.Threshold {
				return true, nil
			}
		}
	default:
		return true, errors.New("Rule couldn't be parsed")
	}
	return false, nil
}

func (a *Alarm) EnableAlarm() {
	a.enabled = true
}

func (a *Alarm) DisableAlarm() {
	a.enabled = false
}

func (a *Alarm) IsEnabled() bool {
	return a.enabled
}
