package main

import (
	// "fmt"
	"github.com/op/go-logging"
	"os"
	"time"
)

var log = logging.MustGetLogger("backend")

var filesettings FileSettings

func main() {
	log.Info("starting Graphite-Monitor")
	filesettings.Filename = "backend.conf"
	filename := "graphmon.log"
	out, err := os.Create(filename)
	defer out.Close()
	if err != nil {
		panic(err)
	}
	backend := logging.NewLogBackend(out, "", 0)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)
	if err = SetupLogging([]logging.Backend{backend, backend2}); err != nil {
		panic(err)
	}
	start()
}

func SetupLogging(backends []logging.Backend) error {
	logging.SetBackend(backends...)
	return nil
}

func start() {
	for {
		log.Debug("Running Loop")
		//check settings
		settings, err := filesettings.UpdateSettings()
		if err != nil {
			log.Warning("Couldn't get settings", err)
			continue
		}
		//generate notifications for alarms
		notifications, err := GenerateNotifications(settings.Alarms, settings.Graphite)
		if err != nil {
			log.Warning("Couldn't Generate Notifications", err)
			continue
		}
		//send notifications
		for _, v := range settings.Notifiers {
			if err = v.Notify(notifications); err != nil {
				log.Warning("Couldn't Notify", err)
			}
		}
		time.Sleep(settings.Frequency)
	}
}

func GenerateNotifications(alarms []Alarm, getter DataGetter) ([]Notification, error) {
	log.Debug("Generating Notifications...")
	notifications := make([]Notification, 0)
	for _, alarm := range alarms {
		data, err := getter.GetDataForAlarm(alarm)
		if err != nil {
			log.Warning("Couldn't get data for Target: %s", alarm.Target)
			continue
		}
		alarmedtargets, err := alarm.Down(data)
		if err != nil {
			log.Warning("Couldn't determine if rule has been met", err)
		}
		for _, target := range alarmedtargets {
			notification := Notification{}
			notification.Message = "Rule: " + alarm.Rule + " has been met for target: " + target
			notifications = append(notifications, notification)
		}
	}
	return notifications, nil
}
