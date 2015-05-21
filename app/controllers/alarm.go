package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lswith/graphite-monitor/app/models"
	"github.com/revel/revel"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Alarms struct {
	*revel.Controller
}

type AlarmState map[models.Target]models.State

func (c Alarms) parseAlarm() (*models.Alarm, error) {
	alarm := new(models.Alarm)
	err := json.NewDecoder(c.Request.Body).Decode(alarm)
	return alarm, err
}

func (c Alarms) CreateAlarm() revel.Result {
	alarm, err := c.parseAlarm()
	if err != nil {
		return c.RenderError(err)
	}
	alarm.Validate(c.Validation)
	if c.Validation.HasErrors() {
		return c.RenderError(errors.New("validation error occured"))
	}
	key, err := AddObject(alarm, AlarmBucket)
	if err != nil {
		return c.RenderError(err)
	}
	return c.RenderText(key)
}

func (c Alarms) ReadAlarms() revel.Result {
	m := make(map[string]*models.Alarm)
	ids, err := GetKeys(AlarmBucket)
	if err != nil {
		return c.RenderError(err)
	}
	for _, id := range ids {
		alarm := new(models.Alarm)
		err = GetObject(id, alarm, AlarmBucket)
		if err != nil {
			return c.RenderError(err)
		}
		m[id] = alarm
	}
	return c.RenderJson(m)
}

func (c Alarms) ReadAlarm(id string) revel.Result {
	alarm := new(models.Alarm)
	err := GetObject(id, alarm, AlarmBucket)
	if err != nil {
		return c.RenderError(err)
	}
	return c.RenderJson(alarm)
}

func (c Alarms) DeleteAlarm(id string) revel.Result {
	err := DeleteObject(id, AlarmBucket)
	if err != nil {
		return c.RenderError(err)
	}
	return c.RenderText("SUCCESS")
}

func (c Alarms) ReadState(id string) revel.Result {
	alarm := new(models.Alarm)
	err := GetObject(id, alarm, AlarmBucket)
	if err != nil {
		return c.RenderError(err)
	}
	state, err := GetState(alarm)
	if err != nil {
		return c.RenderError(err)
	}
	state2 := make(map[models.Target]string)
	for k, v := range state {
		state2[k] = v.String()
	}
	return c.RenderJson(state2)
}

func GetState(a *models.Alarm) (AlarmState, error) {
	state := make(map[models.Target]models.State)
	for _, v := range a.Targets {
		revel.INFO.Printf("getting state for target: %s\n", v)
		targetstate, err := updateFromGraphite(a, v)
		if err != nil {
			revel.ERROR.Println(err)
			return state, err
		}
		state[v] = targetstate
	}
	revel.INFO.Printf("finished getting state\n")
	return state, nil
}

func GetUrlValues(a *models.Alarm, target string) url.Values {
	values := url.Values{}
	values.Set("target", target)
	if a.From != "" {
		values.Add("from", a.From)
	}
	now := time.Now().UTC()
	untilstring := fmt.Sprintf("%.2d:%.2d_%d%.2d%.2d", now.Hour(), now.Minute(), now.Year(), now.Month(), now.Day())
	values.Add("until", untilstring)
	return values

}

func updateFromGraphite(config *models.Alarm, target models.Target) (models.State, error) {
	values := GetUrlValues(config, string(target))
	values.Add("format", "json")
	apiurl := config.Endpoint + "/render?" + values.Encode()
	resp, err := http.Get(apiurl)
	if err != nil {
		revel.ERROR.Println(err)
		return models.NaN, err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	var m []models.Data
	for {
		if err := dec.Decode(&m); err == io.EOF {
			break
		}
		if err != nil {
			revel.ERROR.Println(err)
			return models.NaN, err
		}
	}
	return determineState(m, config.Rule, config.Threshold)
}

func determineState(m []models.Data, rule string, threshold float64) (models.State, error) {
	state := models.GOOD
	for _, d := range m {
		switch rule {
		case "==":
			for j := range d.DataPoints {
				if d.DataPoints[j][0] == threshold {
					state = models.BAD
					break
				}
			}
		case "!=":
			for j := range d.DataPoints {
				if d.DataPoints[j][0] != threshold {
					state = models.BAD
					break
				}
			}
		case "<":
			for j := range d.DataPoints {
				if d.DataPoints[j][0] < threshold {
					state = models.BAD
					break
				}
			}
		case "<=":
			for j := range d.DataPoints {
				if d.DataPoints[j][0] <= threshold {
					state = models.BAD
					break
				}
			}
		case ">":
			for j := range d.DataPoints {
				if d.DataPoints[j][0] > threshold {
					state = models.BAD
					break
				}
			}
		case ">=":
			for j := range d.DataPoints {
				if d.DataPoints[j][0] >= threshold {
					state = models.BAD
					break
				}
			}
		default:
			return models.NaN, fmt.Errorf("rule doesn't match")
		}
	}
	return state, nil
}
