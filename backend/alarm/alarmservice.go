package alarm

import (
	"github.com/lswith/graphite-monitor/backend/storage"
	"log"
)

var bucket string = "alarms"

func CreateAlarm(endpoint string, targets []string, from string, until string, threshold float64, rule string) *Alarm {
	a := newAlarm(endpoint, targets, from, until, threshold, rule)
	return a
}

func StoreAlarm(a *Alarm) (string, error) {
	return storage.AddObject(a, bucket)
}

func DeleteAlarm(id string) error {
	return storage.DeleteObject(id, bucket)
}

func GetAlarm(id string) (*Alarm, error) {
	log.Printf("getting alarm for id: %s\n", id)
	a := new(Alarm)
	err := storage.GetObject(id, a, bucket)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return a, nil
}

func GetAlarms() (map[string]*Alarm, error) {
	keys, err := storage.GetKeys(bucket)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	m := make(map[string]*Alarm)
	for _, v := range keys {
		a := new(Alarm)
		err := storage.GetObject(v, a, bucket)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		m[v] = a
	}
	return m, nil
}
