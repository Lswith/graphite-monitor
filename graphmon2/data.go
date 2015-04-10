package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type DataGetter interface {
	GetDataForTarget(target string, interval string) ([]Data, error)
}

type Data struct {
	Target     string
	DataPoints [][2]float64
}

type GraphiteGetter struct {
	Endpoint string
	Client   *http.Client
}

func (g *GraphiteGetter) GetDataForTarget(target string, interval string) ([]Data, error) {
	values := url.Values{}
	values.Set("target", target)
	values.Add("from", interval)
	values.Add("format", "json")
	u, err := url.ParseRequestURI(g.Endpoint)
	if err != nil {
		log.Error("Couldn't get Data for target: "+target, err)
		return []Data{}, nil
	}
	u.Path = "/render"
	u.RawQuery = values.Encode()
	urlStr := fmt.Sprintf("%v", u)
	resp, err := g.Client.Get(urlStr)
	if err != nil {
		log.Error("Couldn't get Data for target: "+target, err)
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
