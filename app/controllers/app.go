package controllers

import "github.com/revel/revel"
import _ "github.com/lswith/graphite-monitor/app/backend"

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}
