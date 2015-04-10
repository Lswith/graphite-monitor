package main

import (
	"encoding/json"
	"io"
	"net/url"
)

type GetDataFunc func(config Config, client Getter) ([]Data, error)

func GetData(config Config, client Getter) ([]Data, error) {
	ep, _ := url.Parse(config.Endpoint)
	values := url.Values{}
	values.Set("target", config.Target)
	values.Add("from", config.Interval)
	actualurl := ep.String() + "/render" + "?" + values.Encode() + "&format=json"
	resp, err := client.Get(actualurl)
	if err != nil {
		return []Data{}, err
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
