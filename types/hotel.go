package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Hotel struct {
	ID       primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Location string               `bson:"location" json:"location"`
	Name     string               `bson:"name" json:"name"`
	Rating   int                  `bson:"rating" json:"rating"`
	Rooms    []primitive.ObjectID `bson:"rooms" json:"rooms"`
}

type Room struct {
	HotelID primitive.ObjectID `bson:"hotelID" json:"hotelID"`
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Price   float64            `bson:"price" json:"price"`
	Seaside bool               `bson:"seaside" json:"seaside"`
	Size    string             `bson:"size" json:"size"`
}
