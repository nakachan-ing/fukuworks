package test

import (
	"errors"
	"sync"
	"time"

	"github.com/nakachan-ing/fukuworks/backend-go/internal/domain/models"
)

type MockTaskRepo struct {
	tasks  map[string]map[uint][]*models.Task // user -> projectID -> []Task
	taskID uint
	mu     sync.Mutex
}

func NewMockTaskRepo() *MockTaskRepo {
	return &MockTaskRepo{
		tasks:  make(map[string]map[uint][]*models.Task),
		taskID: 1,
	}
}

func (m *MockTaskRepo) Create(userName string, projectID uint, task *models.Task) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	task.ID = m.taskID
	task.ProjectID = projectID
	m.taskID++

	if _, ok := m.tasks[userName]; !ok {
		m.tasks[userName] = make(map[uint][]*models.Task)
	}
	m.tasks[userName][projectID] = append(m.tasks[userName][projectID], task)
	return nil
}

func (m *MockTaskRepo) Find(userName string, projectID uint, taskID uint) (*models.Task, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	tasks, ok := m.tasks[userName][projectID]
	if !ok {
		return nil, errors.New("not authorized or project not found")
	}
	for _, t := range tasks {
		if t.ID == taskID {
			return t, nil
		}
	}
	return nil, errors.New("task not found")
}

func (m *MockTaskRepo) FindAll(userName string, projectID uint) ([]models.Task, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	taskPtrs, ok := m.tasks[userName][projectID]
	if !ok {
		return nil, errors.New("no tasks found")
	}

	tasks := make([]models.Task, len(taskPtrs))
	for i, t := range taskPtrs {
		tasks[i] = *t
	}
	return tasks, nil
}

func (m *MockTaskRepo) Update(userName string, projectID uint, taskID uint, task *models.Task) (*models.Task, error) {
	t, err := m.Find(userName, projectID, taskID)
	if err != nil {
		return nil, err
	}
	t.Title = task.Title
	t.Status = task.Status
	t.DueDate = task.DueDate
	return t, nil
}

func (m *MockTaskRepo) SoftDelete(userName string, projectID uint, taskID uint) error {
	t, err := m.Find(userName, projectID, taskID)
	if err != nil {
		return err
	}
	t.DeletedAt.Valid = true
	t.DeletedAt.Time = time.Now()
	return nil
}

func (m *MockTaskRepo) FindAllTasksForOwner() ([]models.Task, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var allTasks []models.Task
	for _, projects := range m.tasks {
		for _, taskList := range projects {
			for _, t := range taskList {
				allTasks = append(allTasks, *t)
			}
		}
	}
	return allTasks, nil
}

func (m *MockTaskRepo) HardDelete(id uint) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for user, projects := range m.tasks {
		for pid, taskList := range projects {
			for i, t := range taskList {
				if t.ID == id {
					m.tasks[user][pid] = append(taskList[:i], taskList[i+1:]...)
					return nil
				}
			}
		}
	}
	return errors.New("task not found")
}
