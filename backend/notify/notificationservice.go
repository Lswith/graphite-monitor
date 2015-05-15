package notify

import (
	"github.com/lswith/graphite-monitor/backend/storage"
)

var bucket string = "notifiers"

func CreateNotifier() (string, error) {
	n := newNotifier()
	return storage.AddObject(n, bucket)
}

func DeleteNotifier(id string) error {
	return storage.DeleteObject(id, bucket)
}

func GetNotifier(id string) (*Notifier, error) {
	n := new(Notifier)
	err := storage.GetObject(id, n, bucket)
	if err != nil {
		return nil, err
	}
	return n, nil
}

func GetNotifiers() (map[string]*Notifier, error) {
	keys, err := storage.GetKeys(bucket)
	if err != nil {
		return nil, err
	}
	m := make(map[string]*Notifier)
	for _, v := range keys {
		n := new(Notifier)
		err := storage.GetObject(v, n, bucket)
		if err != nil {
			return nil, err
		}
		m[v] = n
	}
	return m, nil
}
