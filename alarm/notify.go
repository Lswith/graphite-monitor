package alarm

import (
	// "net/smtp"
	"fmt"
)

var (
	nm string
)

func notify(state State, a *Alarm) {
	fmt.Printf("Notifier: %s is notifying about Alarm: %s is in State: %t\n", nm, a.Name, state)
}

func init() {
	nm = "test"
}

func UpdateNotifier(name string) {
	nm = name
}
