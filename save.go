package main

import (
	"io"
	"net/http"
)

type Getter interface {
	Get(url string) (resp *http.Response, err error)
}

type SaveFileFunc func(url string, client Getter, file io.Writer) error

func SaveFile(url string, client Getter, file io.Writer) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(file, resp.Body)
	return err
}
