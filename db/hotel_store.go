package db

import (
	"context"

	"github.com/rockchalkwushock/go-hotel-reservations/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelStore interface {
	GetHotelByID(ctx context.Context, id primitive.ObjectID) (*types.Hotel, error)
	GetHotels(ctx context.Context, filter bson.M) ([]*types.Hotel, error)
	InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error)
	UpdateHotel(ctx context.Context, filter bson.M, values bson.M) error
}

type MongoHotelStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client) *MongoHotelStore {
	return &MongoHotelStore{
		client:     client,
		collection: client.Database(DBNAME).Collection("hotels"),
	}
}

func (m *MongoHotelStore) GetHotelByID(ctx context.Context, id primitive.ObjectID) (*types.Hotel, error) {
	var hotel *types.Hotel
	if err := m.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&hotel); err != nil {
		return nil, err
	}
	return hotel, nil
}

func (m *MongoHotelStore) GetHotels(ctx context.Context, filter bson.M) ([]*types.Hotel, error) {
	var hotels []*types.Hotel
	resp, err := m.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if err := resp.All(ctx, &hotels); err != nil {
		return nil, err
	}
	return hotels, nil
}

func (m *MongoHotelStore) InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	resp, err := m.collection.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}

	hotel.ID = resp.InsertedID.(primitive.ObjectID)

	return hotel, nil
}

func (m *MongoHotelStore) UpdateHotel(ctx context.Context, filter bson.M, values bson.M) error {
	_, err := m.collection.UpdateOne(ctx, filter, values)
	if err != nil {
		return err
	}

	return nil
}
