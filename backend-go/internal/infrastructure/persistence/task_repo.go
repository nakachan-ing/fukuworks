package persistence

import (
	"time"

	"github.com/nakachan-ing/fukuworks/backend-go/internal/domain/models"
	"github.com/nakachan-ing/fukuworks/backend-go/internal/domain/repositories"
	"gorm.io/gorm"
)

type TaskRepositoryImpl struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) repositories.TaskRepository {
	return &TaskRepositoryImpl{db}
}

func (r *TaskRepositoryImpl) Create(userName string, pid uint, task *models.Task) error {

	// userIDが欲しい
	var user models.User
	if err := r.db.Where("name = ?", userName).First(&user).Error; err != nil {
		return err
	}

	// projectNameに紐づくprojectIDが欲しい
	var project models.Project
	if err := r.db.Where("user_id = ? AND number= ?", user.ID, pid).First(&project).Error; err != nil {
		return err
	}

	if err := task.SetStatus(task.Status); err != nil {
		return err
	}

	if err := task.SetPriority(task.Priority); err != nil {
		return err
	}

	var maxNumber uint
	r.db.Model(&models.Task{}).
		Where("project_id = ?", project.ID).
		Select("COALESCE(MAX(number), 0)").Scan(&maxNumber)

	task.ProjectID = project.ID
	task.Number = maxNumber + 1

	return r.db.Create(task).Error
}

func (r *TaskRepositoryImpl) Find(userName string, pid, tid uint) (*models.Task, error) {
	var user models.User
	if err := r.db.Where("name = ?", userName).First(&user).Error; err != nil {
		return nil, err
	}

	var project models.Project
	if err := r.db.Where("number = ? AND user_id = ?", pid, user.ID).First(&project).Error; err != nil {
		return nil, err
	}

	var task models.Task
	if err := r.db.Where("number = ? AND project_id = ?", tid, project.ID).First(&task).Error; err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *TaskRepositoryImpl) Update(userName string, pid, tid uint, task *models.Task) (*models.Task, error) {
	existedTask, err := r.Find(userName, pid, tid)
	if err != nil {
		return nil, err
	}

	if err := task.SetStatus(task.Status); err != nil {
		return nil, err
	}

	if err := task.SetPriority(task.Priority); err != nil {
		return nil, err
	}

	updateData := map[string]interface{}{
		"title":       task.Title,
		"description": task.Description,
		"status":      task.Status,
		"priority":    task.Priority,
		"due_date":    task.DueDate,
	}

	return existedTask, r.db.Model(existedTask).Updates(updateData).Error
}

func (r *TaskRepositoryImpl) SoftDelete(userName string, pid, tid uint) error {
	task, err := r.Find(userName, pid, tid)
	if err != nil {
		return err
	}
	return r.db.Model(&models.Task{}).Where("id = ?", task.ID).Update("deleted_at", gorm.DeletedAt{Time: time.Now(), Valid: true}).Error
}

func (r *TaskRepositoryImpl) FindAll(userName string, pid uint) ([]models.Task, error) {
	var user models.User
	if err := r.db.Where("name = ?", userName).First(&user).Error; err != nil {
		return nil, err
	}

	var project models.Project
	if err := r.db.Where("number = ? AND user_id = ?", pid, user.ID).First(&project).Error; err != nil {
		return nil, err
	}

	var tasks []models.Task
	if err := r.db.Where("project_id = ?", project.ID).Find(&tasks).Error; err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *TaskRepositoryImpl) FindAllTasksForOwner() ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.Unscoped().Find(&tasks).Error

	return tasks, err
}

func (r *TaskRepositoryImpl) HardDelete(id uint) error {
	var task models.Task
	if err := r.db.Unscoped().First(&task, id).Error; err != nil {
		return err
	}

	return r.db.Unscoped().Where("id = ?", task.ID).Delete(&models.Task{}).Error
}
