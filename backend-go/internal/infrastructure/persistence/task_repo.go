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
	if err := r.db.First(task, 1).Error; err != nil {
		return err
	}
	return nil
}

func (r *TaskRepositoryImpl) FindByID(id uint) (*models.Task, error) {
	var task models.Task
	if err := r.db.First(task, id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepositoryImpl) FindAll() ([]models.Task, error) {
	var tasks []models.Task
	if err := r.db.Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *TaskRepositoryImpl) Update(task *models.Task) error {
	return nil
}

func (r *TaskRepositoryImpl) Delete(id uint) error {
	return nil
}
