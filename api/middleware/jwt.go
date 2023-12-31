package middleware

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(c *fiber.Ctx) error {
	fmt.Println("--- JWTAuthentication ---")

	token, ok := c.GetReqHeaders()["X-Api-Token"]

	if !ok {
		return fmt.Errorf("Unauthorized")
	}

	fmt.Printf("token: %s\n", token)

	return nil
}

func parseToken(tokenStr string) error {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			fmt.Printf("Unexpected signing method: %v\n", token.Header["alg"])
			return nil, fmt.Errorf("Unauthorized")
		}
		secret := os.Getenv("JWT_SECRET")
		fmt.Println("NEVER PRINT SECRET", secret)
		return []byte(secret), nil
	})

	if err != nil {
		fmt.Printf("Failed to parse token: %v\n", err)
		return fmt.Errorf("Unauthorized")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Printf("claims: %v\n", claims)
	}

	return fmt.Errorf("Unauthorized")
}
