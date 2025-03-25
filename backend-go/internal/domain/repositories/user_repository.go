package repositories

import "github.com/nakachan-ing/fukuworks/backend-go/internal/domain/models"

type UserRepository interface {
	// for user
	Create(user *models.User) error
	Find(userName string) (*models.User, error)
	Update(userName string, user *models.User) (*models.User, error)
	SoftDelete(userName string) error

	// for owner
	FindAll() ([]models.User, error)
	HardDelete(id uint) error
}
