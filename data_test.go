package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetData1(t *testing.T) {
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
						[null, 150]
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
	config := Config{}
	config.Endpoint = "http://example.com"
	config.Target = "examples"
	config.Interval = "-20mins"
	data, err := GetData(config, httpClient)
	if err != nil {
		t.Error("shouldn't have returned an error")
	}
	if len(data) > 1 {
		t.Error("did not parse data correctly")
	}
	datatarget := data[0].Target
	if datatarget != "examples" {
		t.Error("did not parse target correctly")
	}
	datapoints := data[0].DataPoints
	for i, v := range datapoints {
		var x float64 = 100 + (10 * float64(i))
		y := 0.0
		if v[0] != y || v[1] != x {
			t.Errorf("did not parse datapoints correctly:%f,%f\n", v[0], v[1])
		}
	}
}

func TestGetData2(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer server.Close()
	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return &url.URL{}, errors.New("failed request")
		},
	}
	httpClient := &http.Client{Transport: transport}
	config := Config{}
	config.Endpoint = "http://example.com"
	config.Target = "examples"
	config.Interval = "-20mins"
	_, err := GetData(config, httpClient)
	if err == nil {
		t.Error("should have returned an error for a failed request")
	}

}

func TestGetData3(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		jsondata := "hello"
		w.Write([]byte(jsondata))
	}))
	defer server.Close()
	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}
	httpClient := &http.Client{Transport: transport}
	config := Config{}
	config.Endpoint = "http://example.com"
	config.Target = "examples"
	config.Interval = "-20mins"
	_, err := GetData(config, httpClient)
	if err == nil {
		t.Error("should have returned an error for incorrectly parsing bad json")
	}
}
