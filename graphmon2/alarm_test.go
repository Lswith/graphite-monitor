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
			[2]float64{
				1.1,
				2.2,
			},
		},
	},
}
var hasnotbeenmettests = []Test{
	{"==", 0.0},
	{"!=", 1.1},
	{">=", 3.3},
	{"<=", 0.0},
	{"<", 0.0},
	{">", 3.3},
}

var hasbeenmettests = []Test{
	{"==", 1.1},
	{"!=", 0.0},
	{">=", 0.0},
	{"<=", 1.2},
	{"<", 1.2},
	{">", 1.0},
}

func TestHasRuleBeenMet(t *testing.T) {
	alarm := Alarm{}
	for _, v := range hasnotbeenmettests {
		alarm.Rule = v.rule
		alarm.Threshold = v.thres
		alarm.Target = "test"
		targets, err := alarm.HasRuleBeenMet(d)
		if err != nil {
			t.Error("HasRuleBeenMet shouldn't have generated an error")
		}
		if len(targets) > 0 {
			t.Error("rule: " + v.rule + " should not have been met")
		}
	}
	for _, v := range hasbeenmettests {
		alarm.Rule = v.rule
		alarm.Threshold = v.thres

		targets, err := alarm.HasRuleBeenMet(d)
		if err != nil {
			t.Error("HasRuleBeenMet shouldn't have generated an error")
		}
		if len(targets) == 0 {
			t.Error("rule: " + v.rule + " should have been met")
		}
		if len(targets) != 1 {
			t.Error("rule: " + v.rule + "should have been met once")
		}
	}
}
