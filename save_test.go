package main

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

func TestSaveGraph1(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer server.Close()
	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}
	httpClient := &http.Client{Transport: transport}
	url := "http://example.com"
	file, err := os.Create("test.txt")
	defer file.Close()
	if err != nil {
		t.Error(err)
	}
	err = SaveFile(url, httpClient, file)
	if err != nil {
		t.Error("shouldn't have returned an error")
	}
	os.Remove("test.txt")
}

func TestSaveGraph2(t *testing.T) {
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
	url := "http://example.com"
	buf := bytes.NewBuffer(make([]byte, 0))
	err := SaveFile(url, httpClient, buf)
	if err == nil {
		t.Error("an error should have been returned when there is a failed request")
	}
}

type BrokenFile struct{}

func (*BrokenFile) Write(p []byte) (n int, err error) {
	return 0, errors.New("broken file")
}

func TestSaveGraph3(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("hello"))
	}))
	defer server.Close()
	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}
	httpClient := &http.Client{Transport: transport}
	url := "http://example.com"
	buf := BrokenFile{}
	err := SaveFile(url, httpClient, &buf)
	if err == nil {
		t.Error("an error should have been returned when io couldn't copy to a unknown file")
	}
}
