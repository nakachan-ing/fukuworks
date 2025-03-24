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

func (r *ProjectRepositoryImpl) Create(userName string, project *models.Project) error {
	var user models.User
	if err := r.db.Where("name = ?", userName).First(&user).Error; err != nil {
		return err
	}

	if err := project.SetStatus(project.Status); err != nil {
		return err
	}
	var maxNumber uint
	r.db.Model(&models.Project{}).
		Where("user_id = ?", user.ID).
		Select("COALESCE(MAX(number), 0)").Scan(&maxNumber)

	project.UserID = user.ID
	project.Number = maxNumber + 1

	return r.db.Create(project).Error
}

func (r *ProjectRepositoryImpl) Find(userName string, id uint) (*models.Project, error) {
	var user models.User
	if err := r.db.Where("name = ?", userName).First(&user).Error; err != nil {
		return nil, err
	}

	var project models.Project
	if err := r.db.Where("number = ? AND user_id = ?", id, user.ID).First(&project).Error; err != nil {
		return nil, err
	}

	return &project, nil
}

func (r *ProjectRepositoryImpl) Update(userName string, id uint, project *models.Project) (*models.Project, error) {
	existedProject, err := r.Find(userName, id)
	if err != nil {
		return nil, err
	}

	updateData := map[string]interface{}{
		"title":         project.Title,
		"description":   project.Description,
		"platform":      project.Platform,
		"client":        project.Client,
		"estimated_fee": project.EstimatedFee,
		"status":        project.Status,
		"deadline":      project.Deadline,
	}
	return existedProject, r.db.Model(existedProject).Updates(updateData).Error
}

func (r *ProjectRepositoryImpl) SoftDelete(userName string, id uint) error {
	project, err := r.Find(userName, id)
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

func (r *ProjectRepositoryImpl) FindAll(userName string) ([]models.Project, error) {
	var user models.User
	if err := r.db.Where("name = ?", userName).First(&user).Error; err != nil {
		return nil, err
	}

	var projects []models.Project
	if err := r.db.Where("user_id = ?", user.ID).Find(&projects).Error; err != nil {
		return nil, err
	}

	return projects, nil
}

func (r *ProjectRepositoryImpl) FindAllForOwner() ([]models.Project, error) {
	var projects []models.Project
	err := r.db.Unscoped().Find(&projects).Error

	return projects, err
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
