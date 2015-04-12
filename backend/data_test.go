package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetDataForAlarm(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		jsondata := `
			[
				{
					"target": "examples", 
					"datapoints": [
						[0.0, 100], 
						[0.0, 110], 
						[0.0, 120], 
						[0.0, 130], 
						[0.0, 140], 
						[0.0, 150]
					]
				}
			]
		`
		w.Write([]byte(jsondata))
	}))
	defer server.Close()
	var request *http.Request
	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			request = req
			return url.Parse(server.URL)
		},
	}
	httpClient := &http.Client{Transport: transport}
	g := GraphiteGetter{}
	g.Endpoint = "http://test.com"
	g.Client = httpClient
	data, err := g.GetDataForAlarm(Alarm{"stats.*", 0.0, "-5mins", ">", true})
	if err != nil {
		t.Error("GetDataForTarget should not have thrown an error")
	}
	if len(data) != 1 {
		t.Error("GetDataForTarget is not returning the correct amount of datapoints")
	}
	if len(data[0].DataPoints) != 6 {
		t.Error("GetDataForTarget is not returning the correct amount of datapoints")
	}
	for i, v := range data[0].DataPoints {
		if v[0] != 0.0 {
			t.Error("GetDataForTarget has not parsed the datapoints correctly")
		}
		if v[1] != (100.0 + float64(i)*10.0) {
			t.Error("GetDataForTarget has not parsed the datapoints correctly")
		}
	}
	if err = request.ParseForm(); err != nil {
		t.Error("GetDataForTarget created a bad request")
	}
	if request.Host != "test.com" {
		t.Error("GetDataForTarget created a bad request")
	}
	if request.Method != "GET" {
		t.Error("GetDataForTarget created a bad request")
	}
	if request.Form.Get("target") != "stats.*" {
		t.Error("GetDataForTarget created a bad request")
	}
	if request.Form.Get("from") != "-5mins" {
		t.Error("GetDataForTarget created a bad request")
	}
	if request.Form.Get("format") != "json" {
		t.Error("GetDataForTarget created a bad request")
	}
}
