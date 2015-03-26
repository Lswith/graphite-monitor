package main

import (
	"errors"
)

type Alarm struct {
	Target    string
	Rule      string
	Threshold float64
}

type MonitorDataFunc func(d []Data, rule string, thres float64) ([]Alarm, error)

func MonitorData(d []Data, rule string, thres float64) ([]Alarm, error) {
	alarms := make([]Alarm, 0)
	for i := range d {
		data := d[i]
		alarm := Alarm{}
		switch rule {
		case "==":
			for j := range data.DataPoints {
				if data.DataPoints[j][0] == thres {
					alarm.Threshold = thres
					alarm.Target = data.Target
					alarm.Rule = rule
					alarms = append(alarms, alarm)
					break
				}
			}
		case "!=":
			for j := range data.DataPoints {
				if data.DataPoints[j][0] != thres {
					alarm.Threshold = thres
					alarm.Target = data.Target
					alarm.Rule = rule
					alarms = append(alarms, alarm)
					break
				}
			}
		case "<":
			for j := range data.DataPoints {
				if data.DataPoints[j][0] < thres {
					alarm.Threshold = thres
					alarm.Target = data.Target
					alarm.Rule = rule
					alarms = append(alarms, alarm)
					break
				}
			}
		case "<=":
			for j := range data.DataPoints {
				if data.DataPoints[j][0] <= thres {
					alarm.Threshold = thres
					alarm.Target = data.Target
					alarm.Rule = rule
					alarms = append(alarms, alarm)
					break
				}
			}
		case ">":
			for j := range data.DataPoints {
				if data.DataPoints[j][0] > thres {
					alarm.Threshold = thres
					alarm.Target = data.Target
					alarm.Rule = rule
					alarms = append(alarms, alarm)
					break
				}
			}
		case ">=":
			for j := range data.DataPoints {
				if data.DataPoints[j][0] >= thres {
					alarm.Threshold = thres
					alarm.Target = data.Target
					alarm.Rule = rule
					alarms = append(alarms, alarm)
					break
				}
			}
		default:
			return []Alarm{}, errors.New("the rule could not be parsed")
		}
	}

	return alarms, nil
}
