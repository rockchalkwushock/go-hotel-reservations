package db

import (
	"context"

	"github.com/rockchalkwushock/go-hotel-reservations/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomStore interface {
	GetRooms(ctx context.Context, filter bson.M) ([]types.Room, error)
	InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error)
}

type MongoRoomStore struct {
	client     *mongo.Client
	collection *mongo.Collection

	HotelStore
}

func NewMongoRoomStore(client *mongo.Client, hotelStore HotelStore) *MongoRoomStore {
	return &MongoRoomStore{
		client:     client,
		collection: client.Database(DBNAME).Collection("rooms"),
		HotelStore: hotelStore,
	}
}

func (m *MongoRoomStore) GetRooms(ctx context.Context, filter bson.M) ([]types.Room, error) {
	var rooms []types.Room
	resp, err := m.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	if err := resp.All(ctx, &rooms); err != nil {
		return nil, err
	}

	return rooms, nil
}

func (m *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	resp, err := m.collection.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}

	room.ID = resp.InsertedID.(primitive.ObjectID)
	filter := bson.M{"_id": room.HotelID}
	values := bson.M{"$push": bson.M{"rooms": room.ID}}
	if err := m.HotelStore.UpdateHotel(ctx, filter, values); err != nil {
		return nil, err
	}

	return room, nil
}
