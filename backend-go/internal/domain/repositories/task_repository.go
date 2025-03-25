package repositories

import "github.com/nakachan-ing/fukuworks/backend-go/internal/domain/models"

type TaskRepository interface {
	// for user
	Create(userName string, pid uint, task *models.Task) error
	Find(userName string, pid, tid uint) (*models.Task, error)
	Update(userName string, pid, tid uint, task *models.Task) (*models.Task, error)
	SoftDelete(userName string, pid, tid uint) error
	FindAll(userName string, pid uint) ([]models.Task, error)

	// for owner
	FindAllTasksForOwner() ([]models.Task, error)
	HardDelete(id uint) error
}
