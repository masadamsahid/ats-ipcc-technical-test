package middlewares

import (
	"strings"
	"tsilodot/helpers"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
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

	claims, ok := token.Claims.(*helpers.AuthTokenClaims)
	if !ok {
		log.Error().Msg("Invalid token claims structure")
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(&fiber.Map{
			"message": "invalid token claims",
		})
	}

	c.Locals("user", claims)

	return c.Next()
}
