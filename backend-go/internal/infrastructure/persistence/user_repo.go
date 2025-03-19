package persistence

import (
	"github.com/nakachan-ing/fukuworks/backend-go/internal/domain/models"
	"github.com/nakachan-ing/fukuworks/backend-go/internal/domain/repositories"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repositories.UserRepository {
	return &UserRepositoryImpl{db}
}

func (r *UserRepositoryImpl) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepositoryImpl) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *UserRepositoryImpl) FindAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *UserRepositoryImpl) Update(user *models.User) error {
	return nil
}

func (r *UserRepositoryImpl) Delete(id uint) error {
	return nil
}
