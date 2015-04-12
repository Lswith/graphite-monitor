package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type DataGetter interface {
	GetDataForAlarm(alarm Alarm) ([]Data, error)
}

type Data struct {
	Target     string
	DataPoints [][2]float64
}

type GraphiteGetter struct {
	Endpoint string
	Client   *http.Client
}

func (g *GraphiteGetter) GetDataForAlarm(alarm Alarm) ([]Data, error) {
	values := url.Values{}
	values.Add("target", alarm.Target)
	values.Add("from", alarm.Interval)
	values.Add("format", "json")
	u, err := url.ParseRequestURI(g.Endpoint)
	if err != nil {
		log.Error("Couldn't get Data for target: "+alarm.Target, err)
		return []Data{}, nil
	}
	u.Path = "/render"
	u.RawQuery = values.Encode()
	urlStr := u.String()
	resp, err := g.Client.Get(urlStr)
	if err != nil {
		log.Error("Couldn't get Data for target: "+alarm.Target, err)
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	var m []Data
	for {
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			return []Data{}, err
		}
	}
	return m, nil
}
