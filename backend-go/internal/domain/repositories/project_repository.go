package repositories

import "github.com/nakachan-ing/fukuworks/backend-go/internal/domain/models"

type ProjectRepository interface {
	// for user
	Create(project *models.Project) error
	Find(userName string, id uint) (*models.Project, error)
	Update(userName string, projectName string, id uint, project *models.Project) (*models.Project, error)
	SoftDelete(userName string, id uint) error
	FindAll(userName string) ([]models.Project, error)

	// for owner

	HardDelete(id uint) error
}
