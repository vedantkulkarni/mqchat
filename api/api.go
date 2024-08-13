package api

import (
	"fmt"
	"sync"

	"github.com/vedantkulkarni/mqchat/api/routes"
	"github.com/vedantkulkarni/mqchat/pkg/logger"

	"github.com/gofiber/fiber/v3"
)

type API struct {
	HttpAddr     string

}

func NewAPI(httpAddr string) (*API, error) {

	return &API{
		HttpAddr:     httpAddr,
	
	}, nil
}

func (a *API) Start(wg *sync.WaitGroup)  {

	l:= logger.Get()

	defer wg.Done()

	app := fiber.New()
	api := app.Group("/api")

	v1 := api.Group("/v1")

	routes.Init(v1)

	
	err := app.Listen(fmt.Sprintf(":%s", a.HttpAddr))
	if err != nil {
		l.Panic().Err(err).Msg("Error starting API server")	
	}

}
