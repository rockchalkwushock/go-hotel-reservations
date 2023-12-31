package main

import (
	"context"
	"fmt"
	"log"

	"github.com/rockchalkwushock/go-hotel-reservations/db"
	"github.com/rockchalkwushock/go-hotel-reservations/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	ctx        = context.Background()
	hotelStore db.HotelStore
	roomStore  db.RoomStore
)

func init() {
	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
}

func seedHotel(location string, name string, rating int) {
	hotel := types.Hotel{
		Location: location,
		Name:     name,
		Rating:   rating,
		Rooms:    []primitive.ObjectID{},
	}

	rooms := []types.Room{
		{
			Price:   100.00,
			Seaside: false,
			Size:    "small",
		},
		{
			Price:   200.00,
			Seaside: false,
			Size:    "normal",
		},
		{
			Price:   300.00,
			Seaside: false,
			Size:    "king",
		},
		{
			Price:   400.00,
			Seaside: true,
			Size:    "king",
		},
	}
	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted hotel: %+v\n", insertedHotel)
	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		_, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Inserted room: %+v\n", room)
	}
}

func main() {
	seedHotel("San Francisco", "Hotel California", 4)
	seedHotel("New York", "Hilton", 3)
	seedHotel("Miami", "Holiday Inn", 3)
	seedHotel("Los Angeles", "Ritz Carlton", 5)
	seedHotel("San Diego", "Marriott", 5)
}
