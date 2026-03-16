package middlewares

import (
	"log"
	"strings"
	"tsilodot/helpers"

	"github.com/gofiber/fiber/v3"
)

func IsAuthenticated(c fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	if authHeader == "" {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(&fiber.Map{
			"message": "Auth header must be provided",
		})
	}

	// split Bearer and AuthToken
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(&fiber.Map{
			"message": "invalid authorization header format",
		})
	}

	token, err := helpers.VerifyAuthToken(tokenParts[1])
	if err != nil || !token.Valid {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(&fiber.Map{
			"message": "invalid or expired token",
		})
	}

	// log.Printf("%v", token.Claims)

	claims, ok := token.Claims.(*helpers.AuthTokenClaims)
	if !ok {
		log.Println("Error claim:", claims)
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(&fiber.Map{
			"message": "invalid token claims",
		})
	}

	// log.Println("Claims", claims)

	c.Locals("user", claims)

	return c.Next()
}
