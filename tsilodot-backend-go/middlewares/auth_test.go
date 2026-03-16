package middlewares

import (
	"fmt"
	"net/http"
	"testing"
	"tsilodot/helpers"

	"github.com/gofiber/fiber/v3"
)

func TestIsAuthenticated(t *testing.T) {
	app := fiber.New()
	helpers.InitJWT() // ensure secret is set

	app.Get("/protected", IsAuthenticated, func(c fiber.Ctx) error {
		user := c.Locals("user").(*helpers.AuthTokenClaims)
		return c.SendString(fmt.Sprintf("Hello, %s", user.Email))
	})

	t.Run("No Auth Header", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/protected", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("Invalid Format", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "InvalidFormat token")
		resp, _ := app.Test(req)

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("Invalid Token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		resp, _ := app.Test(req)

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("Success", func(t *testing.T) {
		token, err := helpers.CreateAuthToken(helpers.AuthTokenClaims{
			ID:    1,
			Name:  "Test User",
			Email: "test@example.com",
		})
		if err != nil {
			t.Fatalf("Failed to create token: %v", err)
		}

		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		resp, _ := app.Test(req)

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})
}
