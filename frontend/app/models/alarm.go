package models

import (
	"encoding/json"
	"fmt"
	"github.com/revel/revel"
	"io"
	"net/http"
	"net/url"
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

type AlarmState map[Target]State

type State int

type Target string

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

func newAlarm(endpoint string, targets []string, from string, until string, threshold float64, rule string) *Alarm {
	alarmconfig := new(Alarm)
	alarmconfig.Endpoint = endpoint
	alarmconfig.Targets = make([]Target, len(targets))
	for i, v := range targets {
		alarmconfig.Targets[i] = Target(v)
	}
	alarmconfig.From = from
	alarmconfig.Until = until
	alarmconfig.Threshold = threshold
	alarmconfig.Rule = rule
	return alarmconfig
}

func (a *Alarm) GetState() (AlarmState, error) {
	state := make(map[Target]State)
	for _, v := range a.Targets {
		targetstate, err := updateFromGraphite(a, v)
		if err != nil {
			return state, err
		}
		state[v] = targetstate
	}
	return state, nil
}

func updateFromGraphite(config *Alarm, target Target) (State, error) {
	values := url.Values{}
	values.Set("target", string(target))
	if config.From != "" {
		values.Add("from", config.From)
	}
	if config.Until != "" {
		values.Add("until", config.Until)
	}
	values.Add("format", "json")
	apiurl := config.Endpoint + "/render?" + values.Encode()
	resp, err := http.Get(apiurl)
	if err != nil {
		return NaN, err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	var m []data
	for {
		if err := dec.Decode(&m); err == io.EOF {
			break
		}
		if err != nil {
			return NaN, err
		}
	}
	return determineState(m, config.Rule, config.Threshold)
}

func determineState(m []data, rule string, threshold float64) (State, error) {
	state := GOOD
	for _, d := range m {
		switch rule {
		case "==":
			for j := range d.DataPoints {
				if d.DataPoints[j][0] == threshold {
					state = BAD
					break
				}
			}
		case "!=":
			for j := range d.DataPoints {
				if d.DataPoints[j][0] != threshold {
					state = BAD
					break
				}
			}
		case "<":
			for j := range d.DataPoints {
				if d.DataPoints[j][0] < threshold {
					state = BAD
					break
				}
			}
		case "<=":
			for j := range d.DataPoints {
				if d.DataPoints[j][0] <= threshold {
					state = BAD
					break
				}
			}
		case ">":
			for j := range d.DataPoints {
				if d.DataPoints[j][0] > threshold {
					state = BAD
					break
				}
			}
		case ">=":
			for j := range d.DataPoints {
				if d.DataPoints[j][0] >= threshold {
					state = BAD
					break
				}
			}
		default:
			return NaN, fmt.Errorf("rule doesn't match")
		}
	}
	return state, nil
}

type data struct {
	Target     string
	DataPoints [][2]float64
}
