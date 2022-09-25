package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/hadihammurabi/sungoq/constants"
	"github.com/hadihammurabi/sungoq/service"
)

type API struct {
	app *fiber.App

	service *service.Service
	addr    string
}

func New(options ...OptionFunc) (*API, error) {
	api := &API{
		app: fiber.New(),
	}

	for _, opt := range options {
		opt(api)
	}

	if api.addr == "" {
		api.addr = ":8080"
	}

	if api.service == nil {
		return nil, constants.ErrServiceIsEmpty
	}

	return api, nil
}

func (api *API) Route() {
	api.app.Post("/topics", api.CreateTopics)
	api.app.Get("/topics", api.GetAllTopics)
	api.app.Delete("/topics", api.DeleteTopics)

	api.app.Post("/publish", api.Publish)

	api.app.Use("/consume", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	api.app.Get("/consume", websocket.New(api.Consume))
}

func (api *API) Start() {
	api.Route()
	api.app.Listen(api.addr)
}
