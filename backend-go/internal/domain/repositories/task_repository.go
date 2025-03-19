package repositories

import "github.com/nakachan-ing/fukuworks/backend-go/internal/domain/models"

type TaskRepository interface {
	Create(task *models.Task) error
	FindByID(id uint) (*models.Task, error)
	FindAll() ([]models.Task, error)
	Update(task *models.Task) error
	Delete(id uint) error
}
