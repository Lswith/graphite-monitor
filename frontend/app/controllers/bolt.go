package controllers

import (
	"crypto/rand"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/revel/revel"
	"io"
	"log"
)

var (
	Db                    *bolt.DB
	AlarmBucket           string
	NotifierBucket        string
	StatefulWatcherBucket string
	PeriodicWatcherBucket string
)

type BoltController struct {
	*revel.Controller
	Txn *bolt.Tx
}

func (c *BoltController) GenerateKey() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		log.Println(err)
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func (c *BoltController) Begin() revel.Result {
	txn, err := Db.Begin(true)
	if err != nil {
		panic(err)
	}
	c.Txn = txn
	return nil
}

func (c *BoltController) Commit() revel.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Commit(); err != nil {
		panic(err)
	}
	c.Txn = nil
	return nil
}

func (c *BoltController) Rollback() revel.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Rollback(); err != nil {
		panic(err)
	}
	c.Txn = nil
	return nil
}
