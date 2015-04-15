package main

import (
	"testing"
	"time"
)

func TestAlarm(t *testing.T) {
	alm := Alm{"hello"}
	am := new(alarmManager)
	returned := false
	go func() {
		for {
			a, ok := <-am.AddNewAlarm(alm)
			if !ok {
				returned = true
				t.Log("returned set to true")
				break
			}
			t.Logf("recieved a: %s", a)
		}
	}()
	time.Sleep(time.Second * 5)
	if returned == true {
		t.Error("the alarm returned too early")
		t.FailNow()
	}
	am.DeleteAlarm(alm)
	time.Sleep(time.Second * 5)
	if returned == false {
		t.Error("the alarm should have returned now")
		t.FailNow()
	}
}
