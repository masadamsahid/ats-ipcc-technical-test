package controller

import (
	"strconv"
	"tsilodot/dto"
	"tsilodot/helpers"
	"tsilodot/model"
	"tsilodot/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

type TaskController struct {
	taskService service.ITaskService
}

func NewTaskController(taskService service.ITaskService) *TaskController {
	return &TaskController{taskService: taskService}
}

func (t *TaskController) CreateTask(c fiber.Ctx) error {
	userClaims := c.Locals("user").(*helpers.AuthTokenClaims)
	body := new(dto.TaskRequest)

	if err := c.Bind().Body(body); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(&dto.ResponseGeneric[any]{
				Message: "Invalid request body",
				Errors:  err.Error(),
			})
		}

		return c.JSON(&dto.ResponseGeneric[any]{
			Message: "Failed",
			Errors:  helpers.HandleValidationErrors(validationErrors),
		})
	}

	task := &model.Task{
		Title:       body.Title,
		Description: body.Description,
		Status:      body.Status,
		DueDate:     body.GetDueDate(),
	}

	newTask, err := t.taskService.CreateTask(userClaims.ID, task)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(&dto.ResponseGeneric[any]{
			Message: "Failed to create task",
		})
	}

	return c.JSON(&dto.ResponseGeneric[*model.Task]{
		Message: "Task created successfully",
		Data:    &newTask,
	})
}

func (t *TaskController) GetTasks(c fiber.Ctx) error {
	userClaims := c.Locals("user").(*helpers.AuthTokenClaims)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "5"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 5
	}

	offset := (page - 1) * limit

	tasks, total, err := t.taskService.GetTasksByUserID(userClaims.ID, limit, offset)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(&dto.ResponseGeneric[any]{
			Message: "Failed to fetch tasks",
		})
	}

	totalPages := (int(total) + limit - 1) / limit

	return c.JSON(&dto.ResponseGeneric[[]model.Task]{
		Message: "Tasks fetched successfully",
		Pagination: &dto.Pagination{
			CurrentPage: page,
			TotalPages:  totalPages,
			TotalItems:  int(total),
		},
		Data: &tasks,
	})
}

func (t *TaskController) GetTaskByID(c fiber.Ctx) error {
	idParam := c.Params("id")
	taskId, err := strconv.Atoi(idParam)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(&dto.ResponseGeneric[any]{
			Message: "Invalid task ID",
		})
	}

	task, err := t.taskService.GetTaskByID(uint(taskId))
	if err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(&dto.ResponseGeneric[any]{
			Message: "Task not found",
		})
	}

	return c.JSON(&dto.ResponseGeneric[*model.Task]{
		Message: "Task fetched successfully",
		Data:    &task,
	})
}

func (t *TaskController) UpdateTask(c fiber.Ctx) error {
	idParam := c.Params("id")
	taskId, err := strconv.Atoi(idParam)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(&dto.ResponseGeneric[any]{
			Message: "Invalid task ID",
		})
	}

	body := new(dto.TaskRequest)
	if err := c.Bind().Body(body); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(&dto.ResponseGeneric[any]{
				Message: "Invalid request body",
				Errors:  err.Error(),
			})
		}

		return c.JSON(&dto.ResponseGeneric[any]{
			Message: "Failed",
			Errors:  helpers.HandleValidationErrors(validationErrors),
		})
	}

	taskUpdate := &model.Task{
		Title:       body.Title,
		Description: body.Description,
		Status:      body.Status,
		DueDate:     body.GetDueDate(),
	}

	updatedTask, err := t.taskService.UpdateTask(uint(taskId), taskUpdate)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(&dto.ResponseGeneric[any]{
			Message: "Failed to update task",
		})
	}

	return c.JSON(&dto.ResponseGeneric[*model.Task]{
		Message: "Task updated successfully",
		Data:    &updatedTask,
	})
}

func (t *TaskController) DeleteTask(c fiber.Ctx) error {
	idParam := c.Params("id")
	taskId, err := strconv.Atoi(idParam)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(&dto.ResponseGeneric[any]{
			Message: "Invalid task ID",
		})
	}

	err = t.taskService.DeleteTask(uint(taskId))
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(&dto.ResponseGeneric[any]{
			Message: "Failed to delete task",
		})
	}

	return c.JSON(&dto.ResponseGeneric[any]{
		Message: "Task deleted successfully",
	})
}
