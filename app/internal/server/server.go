package server

import (
	"github.com/gofiber/fiber/v2"

	"app/internal/database"

	"github.com/gofiber/template/html/v2"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	engine := html.New("../../templates", ".html")
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "app",
			AppName:      "app",
			Views:        engine,
		}),

		db: database.New(),
	}

	return server
}
