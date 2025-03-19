package persistence

import (
	"github.com/nakachan-ing/fukuworks/backend-go/internal/domain/models"
	"github.com/nakachan-ing/fukuworks/backend-go/internal/domain/repositories"
	"gorm.io/gorm"
)

type ProjectRepositoryImpl struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) repositories.ProjectRepository {
	return &ProjectRepositoryImpl{db}
}

func (r *ProjectRepositoryImpl) Create(project *models.Project) error {
	return r.db.Create(project).Error
}

func (r *ProjectRepositoryImpl) FindByID(id uint) (*models.Project, error) {
	var project models.Project
	err := r.db.First(&project, id).Error
	return &project, err
}

func (r *ProjectRepositoryImpl) FindAll() ([]models.Project, error) {
	var projects []models.Project
	err := r.db.Find(&projects).Error
	return projects, err
}

func (r *ProjectRepositoryImpl) Update(project *models.Project) error {
	return nil
}

func (r *ProjectRepositoryImpl) Delete(id uint) error {
	return nil
}
