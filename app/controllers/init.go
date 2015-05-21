package controllers

import (
	"github.com/boltdb/bolt"
	"github.com/lswith/graphite-monitor/app/models"
	"github.com/revel/revel"
)

func init() {
	revel.OnAppStart(InitDb)
	revel.InterceptFunc(Auth, revel.BEFORE, &Alarms{})
	revel.InterceptFunc(Auth, revel.BEFORE, &Notifiers{})
	revel.InterceptFunc(Auth, revel.BEFORE, &Monitor{})
	revel.InterceptFunc(Auth, revel.BEFORE, &Users{})
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
	defaultusername := getParamString("user.username", "root")
	defaultpassword := getParamString("user.password", "root")
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
	UserBucket = "users"
	PasswordBucket = "passwords"
	err := Db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(AlarmBucket))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(NotifierBucket))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(PeriodicWatcherBucket))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(StatefulWatcherBucket))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(UserBucket))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(PasswordBucket))
		if err != nil {
			return err
		}
		return nil
	})
	user := new(models.User)
	user.Username = defaultusername
	_, err = GetUser(defaultusername)
	if err != nil {
		if err != NoUser {
			revel.ERROR.Fatal(err)
		}
	}
	if err == NoUser {
		id, err := AddObject(user, UserBucket)
		if err != nil {
			revel.ERROR.Fatal(err)
		}
		hashed, err := HashPassword(defaultpassword)
		if err != nil {
			revel.ERROR.Fatal(err)
		}
		err = AddPassword(hashed, id)
		if err != nil {
			revel.ERROR.Fatal(err)
		}
	}
	initMonitor()
}

func initMonitor() {
	Periodicwatchersmap = make(map[string]*models.PeriodicWatcher)
	Statefulwatchersmap = make(map[string]*models.StatefulWatcher)
	RunningWatchersmap = make(map[string]chan bool)
	err := Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(PeriodicWatcherBucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			p := new(models.PeriodicWatcher)
			err := p.UnMarshal(v)
			if err != nil {
				return err
			}
			Periodicwatchersmap[string(k)] = p
			stopchan := make(chan bool)
			RunningWatchersmap[string(k)] = stopchan
			err = RunPeriodicWatcher(string(k))
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		revel.ERROR.Fatal(err)
	}
	err = Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(StatefulWatcherBucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			p := new(models.StatefulWatcher)
			err := p.UnMarshal(v)
			if err != nil {
				return err
			}
			Statefulwatchersmap[string(k)] = p
			stopchan := make(chan bool)
			RunningWatchersmap[string(k)] = stopchan
			err = RunStatefulWatcher(string(k))
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		revel.ERROR.Fatal(err)
	}
}
