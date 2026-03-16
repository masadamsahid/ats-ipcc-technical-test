package repository

import (
	"sync"
	"tsilodot/model"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type TaskRepository struct {
	db *gorm.DB
}

type ITaskRepository interface {
	CreateTask(db *gorm.DB, param *model.Task) (*model.Task, error)
	FindTasksByUserID(db *gorm.DB, userId uint, limit int, offset int) ([]model.Task, int64, error)
	FindTaskByID(db *gorm.DB, taskId uint) (*model.Task, error)
	UpdateTask(db *gorm.DB, task *model.Task) (*model.Task, error)
	DeleteTask(db *gorm.DB, taskId uint) error
}

func NewTaskRepository(db *gorm.DB) ITaskRepository {
	return &TaskRepository{db: db}
}

func (t *TaskRepository) CreateTask(db *gorm.DB, param *model.Task) (*model.Task, error) {
	if db == nil {
		db = t.db
	}
	err := db.Create(param).Error
	if err != nil {
		log.Error().Err(err).Msg("Error creating task")
		return nil, err
	}
	return param, nil
}

func (t *TaskRepository) FindTasksByUserID(db *gorm.DB, userId uint, limit int, offset int) ([]model.Task, int64, error) {
	if db == nil {
		db = t.db
	}
	var tasks []model.Task
	var total int64
	var errCount, errFind error
	var wg sync.WaitGroup // Implementasi concurrent execution (tasks count & tasks list) agar tidak perlu menunggu salah satu request ke DB selesai terlebih dahulu. Karena kedua operasi tidak saling dependen. Tidak ada alasan keduanya harus berjalan sequential dan dimungkinkan eksekusi secara concurrent.

	wg.Add(2)

	go func() { // process 1: Mendapatkan total count dari task dari authenticated user.
		defer wg.Done()
		errCount = db.Model(&model.Task{}).Where("user_id = ?", userId).Count(&total).Error
	}()

	go func() { // process 2: Mendapatkan list task dari authenticated user.
		defer wg.Done()
		errFind = db.Model(&model.Task{}).Where("user_id = ?", userId).Limit(limit).Offset(offset).Find(&tasks).Error
	}()

	wg.Wait()

	if errCount != nil {
		log.Error().Err(errCount).Uint("user_id", userId).Msg("Error counting tasks by user ID")
		return nil, 0, errCount
	}
	if errFind != nil {
		log.Error().Err(errFind).Uint("user_id", userId).Msg("Error finding tasks by user ID")
		return nil, 0, errFind
	}

	return tasks, total, nil
}

func (t *TaskRepository) FindTaskByID(db *gorm.DB, taskId uint) (*model.Task, error) {
	if db == nil {
		db = t.db
	}
	var task model.Task
	err := db.First(&task, taskId).Error
	if err != nil {
		log.Error().Err(err).Uint("task_id", taskId).Msg("Error finding task by ID")
		return nil, err
	}
	return &task, nil
}

func (t *TaskRepository) UpdateTask(db *gorm.DB, task *model.Task) (*model.Task, error) {
	if db == nil {
		db = t.db
	}
	err := db.Save(task).Error
	if err != nil {
		log.Error().Err(err).Uint("task_id", task.ID).Msg("Error updating task")
		return nil, err
	}
	return task, nil
}

func (t *TaskRepository) DeleteTask(db *gorm.DB, taskId uint) error {
	if db == nil {
		db = t.db
	}
	err := db.Delete(&model.Task{}, taskId).Error
	if err != nil {
		log.Error().Err(err).Uint("task_id", taskId).Msg("Error deleting task")
	}
	return err

}
