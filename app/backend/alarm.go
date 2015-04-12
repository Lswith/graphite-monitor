package backend

import (
	"errors"
	// "time"
)

type Alarm struct {
	Target    string
	Threshold float64
	Interval  string
	Rule      string
	Enabled   bool
}

func (a *Alarm) HasRuleBeenMet(datas []Data) ([]string, error) {
	targets := make([]string, 0)
	switch a.Rule {
	case "==":
		for _, data := range datas {
			for j := range data.DataPoints {
				if data.DataPoints[j][0] == a.Threshold {
					targets = append(targets, data.Target)
					break
				}
			}
		}
	case "!=":
		for _, data := range datas {
			for j := range data.DataPoints {
				if data.DataPoints[j][0] != a.Threshold {
					targets = append(targets, data.Target)
					break
				}
			}
		}
	case "<":
		for _, data := range datas {
			for j := range data.DataPoints {
				if data.DataPoints[j][0] < a.Threshold {
					targets = append(targets, data.Target)
					break
				}
			}
		}
	case "<=":
		for _, data := range datas {
			for j := range data.DataPoints {
				if data.DataPoints[j][0] <= a.Threshold {
					targets = append(targets, data.Target)
					break
				}
			}
		}
	case ">":
		for _, data := range datas {
			for j := range data.DataPoints {
				if data.DataPoints[j][0] > a.Threshold {
					targets = append(targets, data.Target)
					break
				}
			}
		}
	case ">=":
		for _, data := range datas {
			for j := range data.DataPoints {
				if data.DataPoints[j][0] >= a.Threshold {
					targets = append(targets, data.Target)
					break
				}
			}
		}
	default:
		return targets, errors.New("Rule couldn't be parsed")
	}

	return targets, nil
}

func (a *Alarm) Down(data []Data) ([]string, error) {
	if a.Enabled {
		return a.HasRuleBeenMet(data)
	}
	return make([]string, 0), nil
}
