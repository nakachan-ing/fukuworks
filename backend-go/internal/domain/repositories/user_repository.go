package repositories

import "github.com/nakachan-ing/fukuworks/backend-go/internal/domain/models"

type UserRepository interface {
	Create(user *models.User) error
	FindByID(id uint) (*models.User, error)
	FindAll() ([]models.User, error)
	Update(id uint, user *models.User) error
	Delete(id uint) error
}
