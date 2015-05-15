package backend

import (
	"github.com/lswith/graphite-monitor/backend/monitor"
	"log"
)

var M *monitor.Monitor

func init() {
	M, err := monitor.CreateMonitor()
	if err != nil {
		log.Fatal(err)
	}
}

//	POST 	/alarms
//	GET		/alarms
//	GET		/alarms/{id}
//	DELETE	/alarms/{id}

//	POST	/notifiers
//	GET		/notifiers
//	GET		/notifiers/{id}
//	DELETE	/notifiers/{id}

//	POST	/watchers
//	GET		/watchers
//	GET		/watchers/{id}
//	DELETE	/watchers/{id}

//	POST	/users
//	GET		/users
//	GET		/users/{id}
//	DELETE	/users/{id}

//Authentication
