package controllers

import (
	"github.com/boltdb/bolt"
	"github.com/lswith/graphite-monitor/frontend/app/models"
	"github.com/revel/revel"
)

func init() {
	revel.OnAppStart(InitDb)
	revel.InterceptMethod((*BoltController).Begin, revel.BEFORE)
	revel.InterceptMethod((*BoltController).Commit, revel.AFTER)
	revel.InterceptMethod((*BoltController).Rollback, revel.FINALLY)
}

func getParamString(param string, defaultValue string) string {
	p, found := revel.Config.String(param)
	if !found {
		if defaultValue == "" {
			revel.ERROR.Fatal("Cound not find parameter: " + param)
		} else {
			return defaultValue
		}
	}
	return p
}

var InitDb func() = func() {
	filename := getParamString("db.filename", "graphite-monitor.db")
	if db, err := bolt.Open(filename, 0600, nil); err != nil {
		revel.ERROR.Fatal(err)
	} else {
		Db = db
	}
	//setup buckets
	AlarmBucket = "alarms"
	NotifierBucket = "notifiers"
	PeriodicWatcherBucket = "periodic"
	StatefulWatcherBucket = "stateful"
	Db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte(AlarmBucket))
		tx.CreateBucketIfNotExists([]byte(NotifierBucket))
		tx.CreateBucketIfNotExists([]byte(PeriodicWatcherBucket))
		tx.CreateBucketIfNotExists([]byte(StatefulWatcherBucket))
		return nil
	})
	initMonitor()
}

func initMonitor() {
	Periodicwatchersmap = make(map[string]*models.PeriodicWatcher)
	Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(PeriodicWatcherBucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			p := new(models.PeriodicWatcher)
			err := p.UnMarshal(v)
			if err != nil {
				return err
			}
			Periodicwatchersmap[string(k)] = p
			go p.Run(Db, AlarmBucket, NotifierBucket)
		}
		return nil
	})
}
