package alarm

import (
	// "net/smtp"
	"fmt"
	"sync"
)

var (
	notifiers     []Notifier
	notifierslock *sync.Mutex
)

type PrintNotifier struct {
}

type Notifier interface {
	Notify(u Update)
}

func (n *PrintNotifier) Notify(u Update) {
	fmt.Printf("Alarm: %s is in State: %t\n", u.A.Name, u.Current)
}

func notify(u Update) {
	notifierslock.Lock()
	for _, n := range notifiers {
		n.Notify(u)
	}
	notifierslock.Unlock()
}

func init() {
	notifiers = make([]Notifier, 0)
	printer := new(PrintNotifier)
	notifiers = append(notifiers, printer)
	notifierslock = new(sync.Mutex)
}
