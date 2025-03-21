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

func (r *UserRepositoryImpl) Find(userName string) (*models.User, error) {
	var user models.User
	err := r.db.Where("name = ?", userName).First(&user).Error
	return &user, err
}

func (r *UserRepositoryImpl) FindAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Unscoped().Find(&users).Error
	return users, err
}

func (r *UserRepositoryImpl) Update(userName string, user *models.User) (*models.User, error) {
	existedUser, err := r.Find(userName)
	if err != nil {
		return nil, err
	}

	updateData := map[string]interface{}{
		"name":  user.Name,
		"email": user.Email,
	}
	return existedUser, r.db.Model(existedUser).Updates(updateData).Error
}

func (r *UserRepositoryImpl) SoftDelete(userName string) error {
	user, err := r.Find(userName)
	if err != nil {
		return err
	}

	if err := r.db.Preload("Projects.Tasks").First(&user, user.ID).Error; err != nil {
		return err
	}

	for _, project := range user.Projects {
		r.db.Model(&models.Task{}).Where("project_id = ?", project.ID).Update("deleted_at", gorm.DeletedAt{Time: time.Now(), Valid: true})
		r.db.Model(&models.Project{}).Where("id = ?", project.ID).Update("deleted_at", gorm.DeletedAt{Time: time.Now(), Valid: true})
	}

	return r.db.Model(&models.User{}).Where("id = ?", user.ID).Update("deleted_at", gorm.DeletedAt{Time: time.Now(), Valid: true}).Error
}

func (r *UserRepositoryImpl) HardDelete(id uint) error {
	var user models.User
	if err := r.db.Unscoped().First(&user, id).Error; err != nil {
		return err
	}

	if err := r.db.Unscoped().Preload("Projects.Tasks").First(&user, id).Error; err != nil {
		return err
	}

	for _, project := range user.Projects {
		r.db.Unscoped().Where("project_id = ?", project.ID).Delete(&models.Task{})
		r.db.Unscoped().Where("id = ?", project.ID).Delete(&models.Project{})
	}

	return r.db.Unscoped().Where("id = ?", id).Delete(&models.User{}).Error
}
