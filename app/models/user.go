package models

import (
	"encoding/json"
	"fmt"
	"github.com/revel/revel"
	"regexp"
)

type User struct {
	Username string
}

type InitUser struct {
	Username string
	Password string
}

func (u *User) String() string {
	return fmt.Sprintf("User(%s)", u.Username)
}

var userRegex = regexp.MustCompile("^\\w*$")

func (user *InitUser) Validate(v *revel.Validation) {
	v.Check(user.Username,
		revel.Required{},
		revel.MaxSize{15},
		revel.MinSize{4},
		revel.Match{userRegex},
	)
}

func (a *User) Marshal() ([]byte, error) {
	return json.Marshal(a)
}

func (a *User) UnMarshal(m []byte) error {
	return json.Unmarshal(m, a)
}
