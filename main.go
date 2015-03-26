package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

type Data struct {
	Target     string
	DataPoints [][2]float64
}

func main() {
	config, err := Setup("graphmon.log", "conf.json")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Graphite-Monitor is starting!")
	Run(config)
}

func Setup(logfile string, configfile string) (Config, error) {
	out, err := os.Create(logfile)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	log.SetOutput(out)
	file, err := os.Open(configfile)
	defer file.Close()
	if err != nil {
		return Config{}, err
	}
	config, err := ReadConfig(file)
	if err != nil {
		return Config{}, err
	}
	return config, err
}

func Run(config Config) {
	defer LogToEmail(config)
	d, err := time.ParseDuration(config.Frequency)
	if err != nil {
		log.Println(err)
	} else {
		d, err = time.ParseDuration("5m")
		if err != nil {
			log.Fatal(err)
		}
	}
	for {
		log.Println("Running Logic")
		Loop(config, GetData, MonitorData, AlarmByEmail)
		time.Sleep(d)
	}
}

func Loop(config Config, getdata GetDataFunc, mondata MonitorDataFunc, alarmbyemail AlarmByEmailFunc) error {
	data, err := getdata(config, &http.Client{})
	if err != nil {
		return err
	}
	alarms, err := mondata(data, config.Rule, config.Threshold)
	if err != nil {
		return err
	}
	for _, alarm := range alarms {
		log.Printf("Target: %s has not met the threshold %f\n", alarm.Target, alarm.Threshold)
		filename := time.Now().Format("01-02-2006T15:04:05") + ".png"
		err := alarmbyemail(alarm, config, filename, SendEmailwithAttachment, SaveFile)
		if err != nil {
			log.Println(err)
			continue
		}
	}
	return nil
}
