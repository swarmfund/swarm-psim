package server

import (
	"sync"

	"gitlab.com/distributed_lab/notificator-server/conf"
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

func (app *App) Start(cfg conf.Config) {
	app.init(cfg)
	app.start(cfg)
}

func (app *App) init(cfg conf.Config) {
	for _, service := range app.services {
		service.Init(cfg)
	}
}

//This method run services in goroutines and will not close while all routines not completed
func (app *App) start(cfg conf.Config) {
	var wg sync.WaitGroup

	wg.Add(len(app.services))
	for _, service := range app.services {
		go app.runService(cfg, &wg, service)
	}

	wg.Wait()
}

func (app *App) runService(cfg conf.Config, wg *sync.WaitGroup, service ServiceInterface) {
	defer wg.Done()
	service.Run(cfg)
}
