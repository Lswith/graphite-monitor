package controllers

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
)

var (
	Db                    *bolt.DB
	AlarmBucket           string
	NotifierBucket        string
	StatefulWatcherBucket string
	PeriodicWatcherBucket string
	UserBucket            string
	PasswordBucket        string
)

type Marshaler interface {
	Marshal() ([]byte, error)
}

type UnMarshaler interface {
	UnMarshal(m []byte) error
}

func generateKey() (string, error) {
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

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func AddPassword(hashedpassword []byte, id string) error {
	return store([]byte(id), hashedpassword, []byte(PasswordBucket))
}

func CheckPassword(password string, id string) error {
	hashedpassword, err := retrieve([]byte(id), []byte(PasswordBucket))
	if err != nil {
		return err
	}
	return bcrypt.CompareHashAndPassword(hashedpassword, []byte(password))
}

func DeletePassword(id string) error {
	return deletefrombolt([]byte(id), []byte(PasswordBucket))
}

func AddObject(m Marshaler, bucket string) (string, error) {
	key, err := generateKey()
	if err != nil {
		log.Println(err)
		return "", err
	}
	value, err := m.Marshal()
	if err != nil {
		log.Println(err)
		return "", err
	}
	return key, store([]byte(key), value, []byte(bucket))
}

func GetObject(key string, u UnMarshaler, bucket string) error {
	log.Printf("retrieving for key: %s\n", key)
	value, err := retrieve([]byte(key), []byte(bucket))
	if err != nil {
		log.Println(err)
		return err
	}
	return u.UnMarshal(value)
}

func DeleteObject(key string, bucket string) error {
	return deletefrombolt([]byte(key), []byte(bucket))
}

func GetKeys(bucket string) ([]string, error) {
	m, err := retrieveAll([]byte(bucket))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	keys := make([]string, 0)
	for k := range m {
		keys = append(keys, string(k))
	}
	return keys, nil
}

func store(key []byte, value []byte, bucket []byte) error {
	err := Db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			log.Println(err)
			return err
		}
		err = b.Put(key, value)
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	})
	return err
}

func retrieve(key []byte, bucket []byte) ([]byte, error) {
	var v []byte
	err := Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		value := b.Get(key)
		if value == nil {
			return errors.New("couldn't find key")
		}
		v = make([]byte, len(value))
		copy(v, value)
		return nil
	})
	return v, err
}

func retrieveAll(bucket []byte) (map[string][]byte, error) {
	data := make(map[string][]byte)
	err := Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			v2 := make([]byte, len(v))
			copy(v2, v)
			data[string(k)] = v2
		}
		return nil
	})
	return data, err
}

func deletefrombolt(key []byte, bucket []byte) error {
	err := Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		err := b.Delete(key)
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	})
	return err
}
