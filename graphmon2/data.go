package main

import (
	"net/url"
)

type DataGetter interface {
	GetDataForTarget(target string) (Data, error)
}

type Data struct {
	DataPoints [][2]float64
}

type GraphiteGetter struct {
	Endpoint url.URL
}

func (g *GraphiteGetter) GetDataForTarget(target string) (Data, error) {
	return Data{}, nil
}
