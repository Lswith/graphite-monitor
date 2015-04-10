package main

import (
	"errors"
	// "time"
)

type Alarm struct {
	Target    string
	Threshold float64
	Rule      string
	Enabled   bool
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

func (a *Alarm) Down(data Data) (bool, error) {
	if a.Enabled {
		return a.HasRuleBeenMet(data)
	}
	return false, nil
}
