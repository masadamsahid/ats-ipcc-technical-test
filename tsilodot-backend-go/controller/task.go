package controller

import (
	"strconv"
	"tsilodot/dto"
	"tsilodot/helpers"
	"tsilodot/model"
	"tsilodot/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
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
		log.Warn().Err(err).Uint("user_id", userClaims.ID).Msg("CreateTask failed: invalid request body")
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

	task := &model.Task{
		Title:       body.Title,
		Description: body.Description,
		Status:      body.Status,
		DueDate:     body.GetDueDate(),
	}

	newTask, err := t.taskService.CreateTask(userClaims.ID, task)
	if err != nil {
		log.Error().Err(err).Uint("user_id", userClaims.ID).Msg("CreateTask failed in service")
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
		log.Error().Err(err).Uint("user_id", userClaims.ID).Msg("GetTasks failed in service")
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
	userClaims := c.Locals("user").(*helpers.AuthTokenClaims)

	idParam := c.Params("id")
	taskId, err := strconv.Atoi(idParam)
	if err != nil {
		log.Warn().Err(err).Uint("user_id", userClaims.ID).Str("task_id_param", idParam).Msg("GetTaskByID failed: invalid task ID")
		c.Status(fiber.StatusBadRequest)
		return c.JSON(&dto.ResponseGeneric[any]{
			Message: "Invalid task ID",
		})
	}

	task, err := t.taskService.GetTaskByID(uint(taskId), userClaims.ID)
	if err != nil {
		log.Warn().Err(err).Uint("user_id", userClaims.ID).Int("task_id", taskId).Msg("GetTaskByID failed in service")
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
	userClaims := c.Locals("user").(*helpers.AuthTokenClaims)

	idParam := c.Params("id")
	taskId, err := strconv.Atoi(idParam)
	if err != nil {
		log.Warn().Err(err).Uint("user_id", userClaims.ID).Str("task_id_param", idParam).Msg("UpdateTask failed: invalid task ID")
		c.Status(fiber.StatusBadRequest)
		return c.JSON(&dto.ResponseGeneric[any]{
			Message: "Invalid task ID",
		})
	}

	body := new(dto.TaskRequest)
	if err := c.Bind().Body(body); err != nil {
		log.Warn().Err(err).Uint("user_id", userClaims.ID).Int("task_id", taskId).Msg("UpdateTask failed: invalid request body")
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

	taskUpdate := &model.Task{
		Title:       body.Title,
		Description: body.Description,
		Status:      body.Status,
		DueDate:     body.GetDueDate(),
	}

	updatedTask, err := t.taskService.UpdateTask(uint(taskId), userClaims.ID, taskUpdate)
	if err != nil {
		log.Error().Err(err).Uint("user_id", userClaims.ID).Int("task_id", taskId).Msg("UpdateTask failed in service")
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(&dto.ResponseGeneric[any]{
			Message: err.Error(),
		})
	}

	return c.JSON(&dto.ResponseGeneric[*model.Task]{
		Message: "Task updated successfully",
		Data:    &updatedTask,
	})
}

func (t *TaskController) DeleteTask(c fiber.Ctx) error {
	userClaims := c.Locals("user").(*helpers.AuthTokenClaims)

	idParam := c.Params("id")
	taskId, err := strconv.Atoi(idParam)
	if err != nil {
		log.Warn().Err(err).Uint("user_id", userClaims.ID).Str("task_id_param", idParam).Msg("DeleteTask failed: invalid task ID")
		c.Status(fiber.StatusBadRequest)
		return c.JSON(&dto.ResponseGeneric[any]{
			Message: "Invalid task ID",
		})
	}

	err = t.taskService.DeleteTask(uint(taskId), userClaims.ID)
	if err != nil {
		log.Error().Err(err).Uint("user_id", userClaims.ID).Int("task_id", taskId).Msg("DeleteTask failed in service")
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(&dto.ResponseGeneric[any]{
			Message: err.Error(),
		})
	}

	return c.JSON(&dto.ResponseGeneric[any]{
		Message: "Task deleted successfully",
	})
}
