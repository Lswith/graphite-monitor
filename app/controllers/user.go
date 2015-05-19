package controllers

import (
	"errors"
	"github.com/boltdb/bolt"
	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	*revel.Controller
}

func Auth(c *revel.Controller) revel.Result {
	if username, password, ok := c.Request.BasicAuth(); ok {
		var hashedpassword []byte = nil
		err := Db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(UserBucket))
			if b == nil {
				return errors.New("bucket doesn't exist")
			}
			p := b.Get([]byte(username))
			if p != nil {
				hashedpassword = make([]byte, len(p))
				copy(hashedpassword, p)
			} else {
				return errors.New("username doesn't exist")
			}
			return nil
		})
		if err != nil {
			revel.ERROR.Println(err)
		}
		err = bcrypt.CompareHashAndPassword(hashedpassword, []byte(password))
		if err == nil {
			return nil
		}
	}
	c.Flash.Error("not authenticated")
	return c.Redirect(App.Index)
}

//must be a unique username
func (u Users) CreateUser() revel.Result {

}

func (u Users) ReadUser(userid string) revel.Result {

}

func (u Users) UpdateUser(userid string) revel.Result {

}

func (u Users) DeleteUser(userid string) revel.Result {

}

func (u Users) CreatePassword(userid string) revel.Result {

}

func (u Users) UpdatePassword(userid string) revel.Result {

}

func (u Users) DeletePassword(userid string) revel.Result {

}

func checkPassword(username string, password string) bool {

}
