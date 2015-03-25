package main

import (
	"log"
)

type Alarm struct {
	Target    string
	Rule      string
	Threshold float64
}

func monitorData(d []Data, rule string, thres float64) []Alarm {
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
			log.Fatal("the rule cannot be parsed!")
		}
	}

	return alarms
}
