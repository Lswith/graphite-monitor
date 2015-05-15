package storage

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"io"
	"log"
)

var filename = "graphite-monitor.db"

var ErrNF = errors.New("couldn't find value for id")

func GenerateKey() (string, error) {
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

type Marshaler interface {
	Marshal() ([]byte, error)
}

type UnMarshaler interface {
	UnMarshal(m []byte) error
}

func AddObject(m Marshaler, bucket string) (string, error) {
	key, err := GenerateKey()
	if err != nil {
		log.Println(err)
		return "", err
	}
	value, err := m.Marshal()
	if err != nil {
		log.Println(err)
		return "", err
	}
	return key, Store([]byte(key), value, []byte(bucket))
}

func GetObject(key string, u UnMarshaler, bucket string) error {
	log.Printf("retrieving for key: %s\n", key)
	value, err := Retrieve([]byte(key), []byte(bucket))
	if err != nil {
		log.Println(err)
		return err
	}
	return u.UnMarshal(value)
}

func DeleteObject(key string, bucket string) error {
	return Delete([]byte(key), []byte(bucket))
}

func GetKeys(bucket string) ([]string, error) {
	m, err := RetrieveAll([]byte(bucket))
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

func Store(key []byte, value []byte, bucket []byte) error {
	db, err := bolt.Open(filename, 0600, nil)
	if err != nil {
		log.Fatal(fmt.Errorf("open bolt %s: %s", filename, err))
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
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

func Retrieve(key []byte, bucket []byte) ([]byte, error) {
	db, err := bolt.Open(filename, 0600, nil)
	if err != nil {
		log.Fatal(fmt.Errorf("open bolt %s: %s", filename, err))
	}
	defer db.Close()
	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists(bucket)
		return nil
	})
	var v []byte
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		value := b.Get(key)
		if value == nil {
			return ErrNF
		}
		v = make([]byte, len(value))
		copy(v, value)
		return nil
	})
	return v, err
}

func RetrieveAll(bucket []byte) (map[string][]byte, error) {
	db, err := bolt.Open(filename, 0600, nil)
	if err != nil {
		log.Fatal(fmt.Errorf("open bolt %s: %s", filename, err))
	}
	defer db.Close()
	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists(bucket)
		return nil
	})
	data := make(map[string][]byte)
	err = db.View(func(tx *bolt.Tx) error {
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

func Delete(key []byte, bucket []byte) error {
	db, err := bolt.Open(filename, 0600, nil)
	if err != nil {
		log.Fatal(fmt.Errorf("open bolt %s: %s", filename, err))
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		err = b.Delete(key)
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	})
	return err
}
