package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetDataForTarget(t *testing.T) {
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
	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}
	httpClient := &http.Client{Transport: transport}
	g := GraphiteGetter{}
	g.Client = httpClient
	data, err := g.GetDataForTarget("target")
	if err != nil {
		t.Error("GetDataForTarget should not have thrown an error")
	}
	if len(data.DataPoints) != 6 {
		t.Error("GetDataForTarget is not returning the correct amount of datapoints")
	}
	for i, v := range data.DataPoints {
		if v[0] != 0.0 {
			t.Error("GetDataForTarget has not parsed the datapoints correctly")
		}
		if v[1] != (100.0 + float64(i)*10.0) {
			t.Error("GetDataForTarget has not parsed the datapoints correctly")
		}
	}
}
