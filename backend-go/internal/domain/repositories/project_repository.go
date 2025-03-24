package repositories

import "github.com/nakachan-ing/fukuworks/backend-go/internal/domain/models"

type ProjectRepository interface {
	// for user
	Create(userName string, project *models.Project) error
	Find(userName string, id uint) (*models.Project, error)
	Update(userName string, id uint, project *models.Project) (*models.Project, error)
	SoftDelete(userName string, id uint) error
	FindAll(userName string) ([]models.Project, error)

	// for owner
	FindAllForOwner() ([]models.Project, error)
	HardDelete(id uint) error
}
