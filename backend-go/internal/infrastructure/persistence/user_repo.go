package persistence

import (
	"time"

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

func (r *UserRepositoryImpl) Update(id uint, user *models.User) error {
	existedUser, err := r.FindByID(id)
	if err != nil {
		return err
	}

	updateData := map[string]interface{}{
		"name":  user.Name,
		"email": user.Email,
	}
	return r.db.Model(existedUser).Updates(updateData).Error
}

func (r *UserRepositoryImpl) Delete(id uint) error {
	var user models.User
	if err := r.db.Preload("Projects.Tasks").First(&user, id).Error; err != nil {
		return err
	}

	for _, project := range user.Projects {
		r.db.Model(&models.Task{}).Where("project_id = ?", project.ID).Update("deleted_at", gorm.DeletedAt{Time: time.Now(), Valid: true})
		r.db.Model(&models.Project{}).Where("id = ?", project.ID).Update("deleted_at", gorm.DeletedAt{Time: time.Now(), Valid: true})
	}

	return r.db.Model(&models.User{}).Where("id = ?", id).Update("deleted_at", gorm.DeletedAt{Time: time.Now(), Valid: true}).Error
}
