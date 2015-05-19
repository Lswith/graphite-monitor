package models

import (
	"encoding/json"
	"github.com/revel/revel"
	"regexp"
)

type Alarm struct {
	Endpoint  string
	Targets   []Target
	From      string
	Until     string
	Threshold float64
	Rule      string
}

func (a *Alarm) Validate(v *revel.Validation) {
	v.Required(a.Endpoint)
	v.Match(a.Endpoint, regexp.MustCompile("^(https?://)?([\\da-z\\.-]+)\\.([a-z\\.]{2,6})([/\\w \\.-]*)*/?$"))
	v.Required(a.Targets)
	v.Required(a.Threshold)
	v.Required(a.Rule)
}

type Target string

type Data struct {
	Target     string
	DataPoints [][2]float64
}

type State int

const (
	NaN State = iota
	GOOD
	BAD
)

func (s State) String() string {
	switch s {
	case NaN:
		return "NaN"
	case GOOD:
		return "GOOD"
	case BAD:
		return "BAD"
	}
	return ""
}

func (a *Alarm) Marshal() ([]byte, error) {
	return json.Marshal(a)
}

func (a *Alarm) UnMarshal(m []byte) error {
	return json.Unmarshal(m, a)
}
