package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"net/url"
	"os"
	"time"
)

type Data struct {
	Target     string
	DataPoints [][2]float64
}

func main() {
	out, err := os.Create("graphmon.log")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	log.SetOutput(out)
	file, err := os.Open("conf.json")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	config := readConfig(file)
	auth := smtp.PlainAuth("", config.EmailUser, config.EmailPassword, config.EmailServer)
	d, err := time.ParseDuration(config.Frequency)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if r := recover(); r != nil {
			log.Println("graphite-monitor encounted an error: ", r)
			sendEmail(config.EmailServer+":"+config.EmailPort, auth, "graphite-monitor encountered an error: "+err.Error(), config.EmailTo, config.EmailFrom)
		}
	}()
	for {

		data := getData(config)
		alarms := monitorData(data, config.Rule, config.Threshold)
		for i := range alarms {
			fmt.Printf("Target: %s has not met the threshold %f\n", alarms[i].Target, alarms[i].Threshold)
			name := saveGraph(alarms[i], config)
			sendEmailwithAttachment(config.EmailServer+":"+config.EmailPort, auth, config.EmailSubject+" "+alarms[i].Target, config.EmailTo, config.EmailFrom, name)
			os.Remove(name)
		}
		time.Sleep(d)
	}
}

func saveGraph(alarm Alarm, config Config) string {
	var graphurl = config.Endpoint + "/render?" + "target=" + alarm.Target + "&from=" + config.Interval
	out, err := os.Create(time.Now().Format("01-02-2015T15.04.05") + ".png")
	if err != nil {
		log.Panic(err)
	}
	defer out.Close()
	resp, err := http.Get(graphurl)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()
	io.Copy(out, resp.Body)
	return out.Name()
}

func getData(config Config) []Data {
	ep, _ := url.Parse(config.Endpoint)
	values := url.Values{}
	values.Set("target", config.Target)
	values.Add("from", config.Interval)
	actualurl := ep.String() + "/render" + "?" + values.Encode() + "&format=json"

	resp, err := http.Get(actualurl)
	defer resp.Body.Close()
	if err != nil {
		log.Panic(err)
	}
	dec := json.NewDecoder(resp.Body)
	var m []Data
	for {
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			log.Panic(err)
		}
	}
	return m
}
