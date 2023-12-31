package main

import (
	"context"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/rockchalkwushock/go-hotel-reservations/api"
	"github.com/rockchalkwushock/go-hotel-reservations/api/middleware"
	"github.com/rockchalkwushock/go-hotel-reservations/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	listenAddr := flag.String("listenAddr", ":5000", "The listen address of the API Server.")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	var (
		// stores
		hotelStore = db.NewMongoHotelStore(client)
		roomStore  = db.NewMongoRoomStore(client, hotelStore)
		userStore  = db.NewMongoUserStore(client)
		store      = &db.Store{
			Hotel: hotelStore,
			Room:  roomStore,
			User:  userStore,
		}

		// handlers
		authHandler  = api.NewAuthHandler(userStore)
		hotelHandler = api.NewHotelHandler(store)
		userHandler  = api.NewUserHandler(store)

		// initialize app
		app = fiber.New(config)

		// route versioning
		auth  = app.Group("/api")
		apiV1 = app.Group("/api/v1", middleware.JWTAuthentication(userStore))
	)

	auth.Post("/auth", authHandler.HandleAuthenticate)

	apiV1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiV1.Get("/user", userHandler.HandleGetUsers)
	apiV1.Get("/user/:id", userHandler.HandleGetUser)
	apiV1.Post("/user", userHandler.HandlePostUser)
	apiV1.Put("/user/:id", userHandler.HandlePutUser)

	apiV1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiV1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiV1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	app.Listen(*listenAddr)
}
