package main

import (
	"testing"
)

type Test struct {
	rule  string
	thres float64
}

var d = []Data{
	Data{
		"hello",
		[][2]float64{
			[2]float64{
				1.1,
				2.2,
			},
		},
	},
}
var noalarmtests = []Test{
	{"==", 0.0},
	{"!=", 1.1},
	{">=", 3.3},
	{"<=", 0.0},
	{"<", 0.0},
	{">", 3.3},
}

var alarmtests = []Test{
	{"==", 1.1},
	{"!=", 0.0},
	{">=", 0.0},
	{"<=", 1.2},
	{"<", 1.2},
	{">", 1.0},
}

func TestMonitorData(t *testing.T) {
	for _, v := range noalarmtests {
		alarms := MonitorData(d, v.rule, v.thres)
		if len(alarms) > 0 {
			t.Error("an alarm should not have been generated")
		}
	}
	for _, v := range alarmtests {
		alarms := MonitorData(d, v.rule, v.thres)
		if len(alarms) == 0 {
			t.Error("an alarm should have been generated: ", v.rule, v.thres)
		}
	}
	recovered := false
	defer func() {
		if r := recover(); r != nil {
			recovered = true
		}
	}()
	MonitorData(d, "=", 0.0)
	if recovered != false {
		t.Error("should have sent a panic when it couldn't parse the rule")
	}
}
