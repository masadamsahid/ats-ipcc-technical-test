package controller

import (
	"tsilodot/dto"
	"tsilodot/helpers"
	"tsilodot/model"
	"tsilodot/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

type AuthController struct {
	authService service.IAuthService
}

func NewAuthController(authService service.IAuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (a *AuthController) Register(c fiber.Ctx) error {
	body := new(dto.RegisterRequest)

	if err := c.Bind().Body(body); err != nil {
		log.Warn().Err(err).Msg("Registration failed: invalid request body")

		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(&dto.ResponseGeneric[any]{
				Message: "Invalid request body",
				Errors:  err.Error(),
			})
		}

		c.Status(fiber.StatusBadRequest)
		return c.JSON(&dto.ResponseGeneric[any]{
			Message: "Failed",
			Errors:  helpers.HandleValidationErrors(validationErrors),
		})
	}

	newUser, err := a.authService.Register(&model.User{
		Name:     body.Name,
		Email:    body.Email,
		Password: body.Password,
	})
	if err != nil {
		log.Error().Err(err).Str("email", body.Email).Msg("Registration failed in service")
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(&dto.ResponseGeneric[any]{
			Message: "Failed registering user",
		})
	}

	return c.JSON(&dto.ResponseGeneric[dto.RegisterResponseData]{
		Message: "Registration successful",
		Data:    newUser,
	})
}

// Login
func (a *AuthController) Login(c fiber.Ctx) error {
	body := new(dto.LoginRequest)

	if err := c.Bind().Body(body); err != nil {
		log.Warn().Err(err).Msg("Login failed: invalid request body")
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(&dto.ResponseGeneric[any]{
				Message: "Invalid request body",
				Errors:  err.Error(),
			})
		}

		c.Status(fiber.StatusBadRequest)
		return c.JSON(&dto.ResponseGeneric[any]{
			Message: "Failed",
			Errors:  helpers.HandleValidationErrors(validationErrors),
		})
	}

	loginData, err := a.authService.Login(body)
	if err != nil {
		log.Warn().Err(err).Str("email", body.Email).Msg("Login failed in service")
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(&dto.ResponseGeneric[any]{
			Message: "Invalid email or password",
		})
	}

	return c.JSON(&dto.ResponseGeneric[dto.LoginResponseData]{
		Message: "Login successful",
		Data:    loginData,
	})
}
