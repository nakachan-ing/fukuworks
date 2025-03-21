package persistence

import (
	"time"

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
	if err := project.SetStatus(project.Status); err != nil {
		return err
	}

	return r.db.Create(project).Error
}

func (r *ProjectRepositoryImpl) Find(projectName string) (*models.Project, error) {
	var project models.Project
	err := r.db.Where("title = ?", projectName).First(&project).Error
	return &project, err
}

func (r *ProjectRepositoryImpl) FindAll() ([]models.Project, error) {
	var projects []models.Project
	err := r.db.Find(&projects).Error
	return projects, err
}

func (r *ProjectRepositoryImpl) Update(projectName string, project *models.Project) (*models.Project, error) {
	// existedProject, err := r.Find(projectName)
	// if err != nil {
	// 	return nil, err
	// }

	// updateData := map[string]interface{}{
	// 	"name":  e.Name,
	// 	"email": user.Email,
	// }
	return nil, nil
}

func (r *ProjectRepositoryImpl) SoftDelete(projectName string) error {
	project, err := r.Find(projectName)
	if err != nil {
		return err
	}

	if err := r.db.Preload("Tasks").First(&project, project.ID).Error; err != nil {
		return err
	}

	for _, task := range project.Tasks {
		r.db.Model(&models.Task{}).Where("project_id = ?", task.ProjectID).Update("deleted_at", gorm.DeletedAt{Time: time.Now(), Valid: true})

	}

	return r.db.Model(&models.Project{}).Where("id = ?", project.ID).Update("deleted_at", gorm.DeletedAt{Time: time.Now(), Valid: true}).Error
}

func (r *ProjectRepositoryImpl) HardDelete(id uint) error {
	var project models.Project
	if err := r.db.Unscoped().First(&project, id).Error; err != nil {
		return err
	}

	if err := r.db.Unscoped().Preload("Tasks").First(&project, project.ID).Error; err != nil {
		return err
	}

	for _, task := range project.Tasks {
		r.db.Unscoped().Where("project_id = ?", task.ProjectID).Delete(&models.Task{})

	}

	return r.db.Unscoped().Where("id = ?", project.ID).Delete(&models.Project{}).Error
}
