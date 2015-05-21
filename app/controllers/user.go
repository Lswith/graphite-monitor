package controllers

import (
	"encoding/json"
	"errors"
	"github.com/lswith/graphite-monitor/app/models"
	"github.com/revel/revel"
)

type Users struct {
	*revel.Controller
}

var (
	NoUser error = errors.New("couldn't find user")
)

func Auth(c *revel.Controller) revel.Result {
	if username, password, ok := c.Request.BasicAuth(); ok {
		id, err := GetUser(username)
		if err == nil {
			err = CheckPassword(password, id)
			if err == nil {
				return nil
			}
			revel.ERROR.Println(err)
		}
		revel.ERROR.Println(err)
	}
	c.Flash.Error("not authenticated")
	return c.Redirect(App.Index)
}

func (c Users) parseInitUser() (*models.InitUser, error) {
	user := new(models.InitUser)
	err := json.NewDecoder(c.Request.Body).Decode(user)
	return user, err
}

//must be a unique username
func (u Users) CreateUser() revel.Result {
	user, err := u.parseInitUser()
	if err != nil {
		return u.RenderError(err)
	}
	user.Validate(u.Validation)
	u.CheckUserisUnique(user.Username)
	if u.Validation.HasErrors() {
		u.Validation.Keep()
		u.FlashParams()
		return u.Redirect(App.Index)
	}
	tmpuser := new(models.User)
	tmpuser.Username = user.Username

	key, err := AddObject(tmpuser, UserBucket)
	if err != nil {
		return u.RenderError(err)
	}
	hashed, err := HashPassword(user.Password)
	if err != nil {
		return u.RenderError(err)
	}
	err = AddPassword(hashed, key)
	if err != nil {
		return u.RenderError(err)
	}
	return u.RenderText(key)
}

func (u Users) CheckUserisUnique(username string) {
	revel.INFO.Printf("checking username: %s\n", username)
	_, err := GetUser(username)
	if err != nil {
		if err != NoUser {
			u.Validation.Error("couldn't check if user exists")
		}
	}
}

func GetUser(username string) (string, error) {
	keys, err := GetKeys(UserBucket)
	if err != nil {
		return "", err
	}
	for _, key := range keys {
		user := new(models.User)
		err := GetObject(key, user, UserBucket)
		if err != nil {
			return "", err
		}
		if username == user.Username {
			return key, nil
		}
	}
	return "", NoUser
}

func (u Users) ReadUser(id string) revel.Result {
	user := new(models.User)
	err := GetObject(id, user, UserBucket)
	if err != nil {
		return u.RenderError(err)
	}
	return u.RenderJson(user)
}

func (u Users) ReadUsers() revel.Result {
	m := make(map[string]*models.User)
	ids, err := GetKeys(UserBucket)
	if err != nil {
		return u.RenderError(err)
	}
	for _, id := range ids {
		user := new(models.User)
		err = GetObject(id, user, UserBucket)
		if err != nil {
			return u.RenderError(err)
		}
		m[id] = user
	}
	return u.RenderJson(m)
}

func (u Users) DeleteUser(id string) revel.Result {
	err := DeleteObject(id, UserBucket)
	if err != nil {
		return u.RenderError(err)
	}
	err = DeletePassword(id)
	if err != nil {
		return u.RenderError(err)
	}
	return u.RenderText("SUCCESS")
}

func (c Users) parsePassword() (*models.Password, error) {
	password := new(models.Password)
	err := json.NewDecoder(c.Request.Body).Decode(password)
	return password, err
}

func (u Users) UpdatePassword(id string) revel.Result {
	password, err := u.parsePassword()
	if err != nil {
		return u.RenderError(err)
	}
	err = CheckPassword(password.Oldpassword, id)
	if err != nil {
		return u.RenderError(err)
	}
	hashedpassword, err := HashPassword(password.Newpassword)
	if err != nil {
		return u.RenderError(err)
	}
	err = AddPassword(hashedpassword, id)
	if err != nil {
		return u.RenderError(err)
	}
	return u.RenderText("SUCCESS")
}
