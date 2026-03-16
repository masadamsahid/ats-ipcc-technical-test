package routes

import (
	"tsilodot/controller"
	"tsilodot/middlewares"

	"github.com/gofiber/fiber/v3"
)

func SetupTaskRoutes(app fiber.Router, taskController *controller.TaskController) {
	taskRoute := app.Group("/tasks", middlewares.IsAuthenticated)

	taskRoute.Post("/", taskController.CreateTask)
	taskRoute.Get("/", taskController.GetTasks)
	taskRoute.Get("/:id", taskController.GetTaskByID)
	taskRoute.Put("/:id", taskController.UpdateTask)
	taskRoute.Delete("/:id", taskController.DeleteTask)
}
