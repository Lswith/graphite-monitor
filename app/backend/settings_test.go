package backend

import (
	"os"
	"testing"
	"time"
)

var testsettings = `
	{
		"endpoint":"testendpoint",
		"frequency":"5m",
		"alarms":[
			{
				"target":"stats.*.reads",
				"threshold":0.0,
				"interval":"-20mins",
				"rule":">",
				"enabled":false
			}
		]
	}
`

func TestFileUpdateSettings(t *testing.T) {
	out, err := os.Create("test.conf")
	if err != nil {
		t.Error(err)
	}
	defer out.Close()
	defer os.Remove("test.conf")
	out.WriteString(testsettings)
	files := FileSettings{"test.conf"}
	settings, err := files.UpdateSettings()
	if err != nil {
		t.Error(err)
	}
	if settings.Graphite.Endpoint != "testendpoint" {
		t.Error("FileSettings.UpdateSettings didn't update the graphite Endpoint correctly")
	}
	if settings.Graphite.Client == nil {
		t.Error("FileSettings.UpdateSettings didn't create a correct graphite client")
	}
	if settings.Frequency != time.Minute*5 {
		t.Error("FileSettings.UpdateSettings didn't update the frequency correctly")
	}
	if len(settings.Alarms) != 1 {
		t.Fatal("FileSettings.UpdateSettings didn't update the alarms correctly")
	}
	if settings.Alarms[0].Target != "stats.*.reads" {
		t.Error("FileSettings.UpdateSettings didn't update the alarms correctly")
	}
	if settings.Alarms[0].Threshold != 0.0 {
		t.Error("FileSettings.UpdateSettings didn't update the alarms correctly")
	}
	if settings.Alarms[0].Rule != ">" {
		t.Error("FileSettings.UpdateSettings didn't update the alarms correctly")
	}
	if settings.Alarms[0].Interval != "-20mins" {
		t.Error("FileSettings.UpdateSettings didn't update the alarms correctly")
	}
	if settings.Alarms[0].Enabled != false {
		t.Error("FileSettings.UpdateSettings didn't update the alarms correctly")
	}
}
