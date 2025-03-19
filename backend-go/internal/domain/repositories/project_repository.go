package repositories

import "github.com/nakachan-ing/fukuworks/backend-go/internal/domain/models"

type ProjectRepository interface {
	Create(project *models.Project) error
	FindByID(id uint) (*models.Project, error)
	FindAll() ([]models.Project, error)
	Update(project *models.Project) error
	Delete(id uint) error
}
