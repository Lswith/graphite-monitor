package backend

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
)

type SettingsUpdater interface {
	UpdateSettings() (Settings, error)
}

type Settings struct {
	Frequency time.Duration
	Alarms    []Alarm
	Notifiers []Notifier
	Graphite  *GraphiteGetter
}

type FileSettings struct {
	Filename string
}

type Config struct {
	Alarms    []Alarm
	Endpoint  string
	Frequency string
}

func (fs *FileSettings) UpdateSettings() (Settings, error) {
	var settings = Settings{}
	file, err := os.Open(fs.Filename)
	if err != nil {
		return Settings{}, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Config{}
	err = decoder.Decode(&configuration)
	if err != nil {
		return Settings{}, err
	}
	settings.Frequency, err = time.ParseDuration(configuration.Frequency)
	if err != nil {
		return Settings{}, err
	}
	settings.Graphite = new(GraphiteGetter)
	settings.Graphite.Client = new(http.Client)
	settings.Graphite.Endpoint = configuration.Endpoint
	notifiers := make([]ConsoleNotifier, 1)
	settings.Notifiers = []Notifier{&notifiers[0]}
	settings.Alarms = configuration.Alarms
	return settings, nil
}
