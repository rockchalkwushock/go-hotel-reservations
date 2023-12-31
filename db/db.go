package db

const (
	DBNAME     = "hotel-reservations"
	DBURI      = "mongodb://localhost:27017"
	TESTDBNAME = "hotel-reservations-test"
)

type Store struct {
	Hotel HotelStore
	Room  RoomStore
	User  UserStore
}
