package main

import (
	// "fmt"
	"github.com/op/go-logging"
	"os"
	"time"
)

var log = logging.MustGetLogger("graphmon")

var settings *Settings = new(Settings)

func main() {
	filename := "graphmon.log"
	out, err := os.Create(filename)
	defer out.Close()
	if err != nil {
		panic(err)
	}
	backend := logging.NewLogBackend(out, "", 0)
	if err = SetupLogging(backend); err != nil {
		panic(err)
	}
	settings.UpdateSettings()
	start()
}

func SetupLogging(backend *logging.LogBackend) error {
	logging.SetBackend(backend)
	return nil
}

func start() {
	for {
		//check settings
		settings.UpdateSettings()
		//generate notifications for alarms
		notifications, err := GenerateNotifications(settings.Alarms, settings.Graphite)
		if err != nil {
			log.Error("Couldn't Generate Notifications", err)
			continue
		}
		//send notifications
		for _, v := range settings.Notifiers {
			if err = v.Notify(notifications); err != nil {
				log.Error("Couldn't Notify", err)
			}
		}
		time.Sleep(settings.Frequency)
	}
}

func GenerateNotifications(alarms []Alarm, getter DataGetter) ([]Notification, error) {
	notifications := make([]Notification, 0)
	for _, v := range alarms {
		d, err := getter.GetDataForTarget(v.Target)
		if err != nil {
			log.Error("Couldn't get data for Target: %s", v.Target)
			continue
		}
		down, err := v.Down(d)
		if err != nil {
			log.Error("Couldn't determine if rule has been met", err)
		}
		if down {
			notification := Notification{}
			notification.Message = "Rule: " + v.Rule + " has been met for target: " + v.Target
			notifications = append(notifications, notification)
		}
	}
	return notifications, nil
}
