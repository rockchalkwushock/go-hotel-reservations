package types

import (
	"fmt"
	"regexp"

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

type User struct {
	// `json:"id,omitempty"` will remove the field from the JSON output if the field is empty.
	// `json:"-"` will remove the field completely from the JSON output regardless of its contents.
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"encryptedPassword" json:"-"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	LastName          string             `bson:"lastName" json:"lastName"`
}

func (params CreateUserParams) Validate() []string {
	errors := []string{}
	if len(params.FirstName) < minFirstNameLength {
		errors = append(errors, fmt.Sprintf("first name must be at least %d characters long", minFirstNameLength))
	}
	if len(params.LastName) < minLastNameLength {
		errors = append(errors, fmt.Sprintf("last name must be at least %d characters long", minLastNameLength))
	}
	if len(params.Password) < minPasswordLength {
		errors = append(errors, fmt.Sprintf("password must be at least %d characters long", minPasswordLength))
	}
	if !isEmailValid(params.Email) {
		errors = append(errors, "email is invalid")
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
