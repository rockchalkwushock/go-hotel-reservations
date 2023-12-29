package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rockchalkwushock/go-hotel-reservations/types"
)

func HandleGetUsers(c *fiber.Ctx) error {
	u := types.User{
		FirstName: "James",
		ID:        "",
		LastName:  "Foo",
	}
	return c.JSON(u)
}

func HandleGetUser(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"user": "James Foo"})
}
