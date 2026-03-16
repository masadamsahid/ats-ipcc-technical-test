package main

import (
	"fmt"
	"os"
	"strconv"
	"tsilodot/controller"
	"tsilodot/db"
	"tsilodot/helpers"
	"tsilodot/repository"
	"tsilodot/routes"
	"tsilodot/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"

	scalargo "github.com/bdpiprava/scalar-go"
	swagger "github.com/gofiber/contrib/v3/swaggerui"
)

func init() {
	helpers.InitLogger()
	if err := godotenv.Load(); err != nil {
		log.Info().Msg("No .env file found, using environment variables from OS")
	} else {
		log.Info().Msg("Loading .env success")
	}
}

type structValidator struct {
	validate *validator.Validate
}

// Validator needs to implement the Validate method
func (v *structValidator) Validate(out any) error {
	return v.validate.Struct(out)
}

func main() {
	app := fiber.New(fiber.Config{
		StructValidator: &structValidator{
			validate: validator.New(),
		},
	})
	APP_PORT, err := strconv.Atoi(os.Getenv("APP_PORT"))
	if err != nil {
		APP_PORT = 8080
	}

	helpers.InitJWT()

	defer db.StopDBConnection()
	db.InitDBConnection()

	db.InitRedisConnection()
	defer db.StopRedisConnection()

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// Raw OpenAPI YAML file
	openapiYAMLDocs, err := os.ReadFile("./docs/openapi.yaml")
	if err != nil {
		log.Error().Err(err).Msg("Error loading OpenAPI .yaml file")
	}
	openapiStr := string(openapiYAMLDocs)
	app.Get("/openapi", func(c fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, "text/openapi+yaml; charset=utf-8")
		return c.SendString(string(openapiStr))
	})

	// Swagger
	swaggerCfg := swagger.Config{
		BasePath: "/", // swagger ui base path
		FilePath: "./docs/openapi.yaml",
		Path:     "/swagger",
	}
	app.Use(swagger.New(swaggerCfg))

	// Scalar UI
	html, err := scalargo.NewV2(
		scalargo.WithSpecDir("./docs"),
		scalargo.WithBaseFileName("openapi.yaml"),
		scalargo.WithTheme(scalargo.ThemeBluePlanet),
	)
	if err != nil {
		log.Error().Err(err).Msg("Error Scalar-Go")
	}
	// log.Debug().Str("html", html).Msg("Scalar HTML generated")
	app.Get("/scalar", func(c fiber.Ctx) error {

		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		return c.SendString(html)
	})

	// APIs
	api := app.Group("/api")

	// /api/auth
	userRepo := repository.NewUserRepository(db.DB)
	authService := service.NewAuthService(userRepo)
	authController := controller.NewAuthController(authService)
	routes.SetupAuthRoutes(api, authController)

	// /api/tasks
	taskRepo := repository.NewTaskRepository(db.DB)
	taskService := service.NewTaskService(taskRepo, db.RedisClient)
	taskController := controller.NewTaskController(taskService)
	routes.SetupTaskRoutes(api, taskController)

	log.Info().Msg("Starting server...")
	for _, route := range app.GetRoutes() {
		log.Debug().Str("method", route.Method).Str("path", route.Path).Msg("Available Route")
	}

	log.Fatal().Err(app.Listen(fmt.Sprintf(":%d", APP_PORT))).Msg("Server failed to listen")

}
