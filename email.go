package main

import (
	"github.com/scorredoira/email"
	"log"
	"net/http"
	"net/smtp"
	"os"
)

type EmailerFunc func(string, smtp.Auth, *email.Message) error

type SendEmailwithAttachmentFunc func(addr string, auth smtp.Auth, subject string, to string, from string, filename string) error

func SendEmailwithAttachment(addr string, auth smtp.Auth, subject string, to string, from string, filename string) error {
	return SendEmailwithAttachmentwithEmailer(addr, auth, subject, to, from, filename, email.Send)
}

func SendEmailwithAttachmentwithEmailer(addr string, auth smtp.Auth, subject string, to string, from string, filename string, e EmailerFunc) error {
	m := email.NewMessage(subject, "")
	m.To = []string{to}
	m.From = from
	err := m.Attach(filename)
	if err != nil {
		return err
	}
	err = e(addr, auth, m)
	if err != nil {
		return err
	}
	return nil
}

type SendEmailFunc func(addr string, auth smtp.Auth, subject string, to string, from string) error

func SendEmail(addr string, auth smtp.Auth, subject string, to string, from string) error {
	return SendEmailwithEmailer(addr, auth, subject, to, from, email.Send)
}

func SendEmailwithEmailer(addr string, auth smtp.Auth, subject string, to string, from string, e EmailerFunc) error {
	m := email.NewMessage(subject, "")
	m.To = []string{to}
	m.From = from
	err := e(addr, auth, m)
	if err != nil {
		return err
	}
	return nil
}

type AlarmByEmailFunc func(alarm Alarm, config Config, graphname string, emailsend SendEmailwithAttachmentFunc, save SaveFileFunc) error

func AlarmByEmail(alarm Alarm, config Config, filename string, emailsend SendEmailwithAttachmentFunc, save SaveFileFunc) error {
	auth := smtp.PlainAuth("", config.EmailUser, config.EmailPassword, config.EmailServer)
	client := http.Client{}
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer os.Remove(out.Name())
	defer out.Close()
	var graphurl = config.Endpoint + "/render?" + "target=" + alarm.Target + "&from=" + config.Interval
	err = save(graphurl, &client, out)
	if err != nil {
		return err
	}
	err = emailsend(config.EmailServer+":"+config.EmailPort, auth, config.EmailSubject+" "+alarm.Target, config.EmailTo, config.EmailFrom, out.Name())
	if err != nil {
		return err
	}
	return err
}

func LogToEmail(config Config) {
	auth := smtp.PlainAuth("", config.EmailUser, config.EmailPassword, config.EmailServer)
	if r := recover(); r != nil {
		log.Println("graphite-monitor encounted an error: ", r)
		err := SendEmail(config.EmailServer+":"+config.EmailPort,
			auth,
			"graphite-monitor encountered an error",
			config.EmailTo,
			config.EmailFrom)
		if err != nil {
			log.Println(err)
		}
	}
}
