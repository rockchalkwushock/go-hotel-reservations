package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost         = 12
	minFirstNameLength = 2
	minLastNameLength  = 2
	minPasswordLength  = 7
)

type CreateUserParams struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
}

type UpdateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (p UpdateUserParams) ToBSON() bson.M {
	m := bson.M{}
	if len(p.FirstName) > 0 {
		m["firstName"] = p.FirstName
	}
	if len(p.LastName) > 0 {
		m["lastName"] = p.LastName
	}
	return m
}

type User struct {
	// `json:"id,omitempty"` will remove the field from the JSON output if the field is empty.
	// `json:"-"` will remove the field completely from the JSON output regardless of its contents.
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"encryptedPassword" json:"-"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	LastName          string             `bson:"lastName" json:"lastName"`
}

func (params CreateUserParams) Validate() map[string]string {
	errors := map[string]string{}
	if len(params.FirstName) < minFirstNameLength {
		errors["firstName"] = fmt.Sprintf("first name must be at least %d characters long", minFirstNameLength)
	}
	if len(params.LastName) < minLastNameLength {
		errors["lastName"] = fmt.Sprintf("last name must be at least %d characters long", minLastNameLength)
	}
	if len(params.Password) < minPasswordLength {
		errors["password"] = fmt.Sprintf("password must be at least %d characters long", minPasswordLength)
	}
	if !isEmailValid(params.Email) {
		errors["email"] = "email is invalid"
	}
	return errors
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	return &User{
		Email:             params.Email,
		EncryptedPassword: string(encpw),
		FirstName:         params.FirstName,
		LastName:          params.LastName,
	}, nil
}
