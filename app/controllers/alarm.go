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

func (c Alarms) Add() revel.Result {
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

func (c Alarms) Map() revel.Result {
	m := make(map[string]*models.Alarm)
	//TODO: make map
	return c.RenderJson(m)
}

func (c Alarms) Get(id string) revel.Result {
	alarm := new(models.Alarm)
	err := GetObject(id, alarm, AlarmBucket)
	if err != nil {
		return c.RenderError(err)
	}
	return c.RenderJson(alarm)
}

func (c Alarms) Delete(id string) revel.Result {
	err := DeleteObject(id, AlarmBucket)
	if err != nil {
		return c.RenderError(err)
	}
	return c.RenderText("SUCCESS")
}

func GetState(a *models.Alarm) (AlarmState, error) {
	state := make(map[models.Target]models.State)
	for _, v := range a.Targets {
		targetstate, err := updateFromGraphite(a, v)
		if err != nil {
			return state, err
		}
		state[v] = targetstate
	}
	return state, nil
}

func updateFromGraphite(config *models.Alarm, target models.Target) (models.State, error) {
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
