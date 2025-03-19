package persistence

import (
	"github.com/nakachan-ing/fukuworks/backend-go/internal/domain/models"
	"github.com/nakachan-ing/fukuworks/backend-go/internal/domain/repositories"
	"gorm.io/gorm"
)

type TaskRepositoryImpl struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) repositories.TaskRepository {
	return &TaskRepositoryImpl{db}
}

func (r *TaskRepositoryImpl) Create(task *models.Task) error {
	return r.db.Create(task).Error
}

func (r *TaskRepositoryImpl) FindByID(id uint) (*models.Task, error) {
	var task models.Task
	err := r.db.First(task, id).Error
	return &task, err
}

func (r *TaskRepositoryImpl) FindAll() ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.Find(&tasks).Error
	return tasks, err
}

func (r *TaskRepositoryImpl) Update(task *models.Task) error {
	return nil
}

func (r *TaskRepositoryImpl) Delete(id uint) error {
	return nil
}
