package server

import (
	"sync"
)

type App struct {
	services []ServiceInterface
}

func NewApp() *App {
	return &App{
		services: []ServiceInterface{},
	}
}

func (app *App) AddService(service ServiceInterface) {
	app.services = append(app.services, service)
}

func (app *App) Start() {
	app.init()
	app.start()
}

func (app *App) init() {
	for _, service := range app.services {
		service.Init()
	}
}

func (app *App) start() {
	var wg sync.WaitGroup

	wg.Add(len(app.services))

	for _, service := range app.services {
		go app.runService(&wg, service)
	}

	wg.Wait()
}

func (app *App) runService(wg *sync.WaitGroup, service ServiceInterface) {
	defer wg.Done()
	service.Run()
}
