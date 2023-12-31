package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/rockchalkwushock/go-hotel-reservations/db"
	"github.com/rockchalkwushock/go-hotel-reservations/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName    = "hotel-reservation-test"
	testDbUri = "mongodb://localhost:27017"
)

type testDb struct {
	store *db.Store
}

func setup(t *testing.T) *testDb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testDbUri))
	if err != nil {
		t.Fatalf("Error connecting to mongo: %v", err)
	}

	store := &db.Store{
		User: db.NewMongoUserStore(client),
	}

	return &testDb{
		store: store,
	}
}

func (tbd *testDb) teardown(t *testing.T) {
	if err := tbd.store.User.Drop(context.TODO()); err != nil {
		t.Fatalf("Error dropping db: %v", err)
	}
}

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.store)
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		Email:     "foo@bar.com",
		FirstName: "Foo",
		LastName:  "Bar",
		Password:  "foobar123",
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	var user types.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		t.Error(err)
	}
	if user.Email != params.Email {
		t.Errorf("Expected email to be %s, got %s", params.Email, user.Email)
	}
	if len(user.EncryptedPassword) > 0 {
		t.Errorf("Expected encrypted password to be empty, got %s", user.EncryptedPassword)
	}
	if user.FirstName != params.FirstName {
		t.Errorf("Expected first name to be %s, got %s", params.FirstName, user.FirstName)
	}
	if len(user.ID) == 0 {
		t.Errorf("Expected user to have an ID")
	}
	if user.LastName != params.LastName {
		t.Errorf("Expected last name to be %s, got %s", params.LastName, user.LastName)
	}
}
