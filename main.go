package main

import (
	"flag"

	"github.com/gofiber/fiber/v2"
	"github.com/rockchalkwushock/go-hotel-reservations/api"
)

func main() {
	listenAddr := flag.String("listenAddr", ":5000", "The listen address of the API Server.")
	flag.Parse()

	app := fiber.New()
	apiV1 := app.Group("/api/v1")

	apiV1.Get("/user", api.HandleGetUsers)
	apiV1.Get("/user/:id", api.HandleGetUser)

	app.Listen(*listenAddr)
}
