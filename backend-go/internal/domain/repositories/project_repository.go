package repositories

import "github.com/nakachan-ing/fukuworks/backend-go/internal/domain/models"

type ProjectRepository interface {
	// for user
	Create(project *models.Project) error
	Find(projectName string) (*models.Project, error)
	Update(projectName string, project *models.Project) (*models.Project, error)
	SoftDelete(projectName string) error

	// for owner
	FindAll() ([]models.Project, error)
	HardDelete(id uint) error
}
