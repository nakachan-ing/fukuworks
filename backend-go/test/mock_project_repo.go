package test

import (
	"errors"
	"sync"
	"time"

	"github.com/nakachan-ing/fukuworks/backend-go/internal/domain/models"
	"github.com/nakachan-ing/fukuworks/backend-go/internal/domain/repositories"
)

type MockProjectRepo struct {
	mu        sync.Mutex
	projects  map[uint]*models.Project
	userIndex map[string][]*models.Project
	idSeq     uint
}

func NewMockProjectRepo() repositories.ProjectRepository {
	return &MockProjectRepo{
		projects:  make(map[uint]*models.Project),
		userIndex: make(map[string][]*models.Project),
		idSeq:     1,
	}
}

func (m *MockProjectRepo) Create(userName string, project *models.Project) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	project.ID = m.idSeq
	project.Number = m.idSeq
	project.CreatedAt = time.Now()
	project.UpdatedAt = time.Now()
	m.projects[project.ID] = project
	m.userIndex[userName] = append(m.userIndex[userName], project)
	m.idSeq++
	return nil
}

func (m *MockProjectRepo) FindAll(userName string) ([]models.Project, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	list := []models.Project{}
	for _, p := range m.userIndex[userName] {
		list = append(list, *p)
	}
	return list, nil
}

func (m *MockProjectRepo) Find(userName string, id uint) (*models.Project, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	p, ok := m.projects[id]
	if !ok {
		return nil, errors.New("project not found")
	}

	for _, p2 := range m.userIndex[userName] {
		if p2.ID == id {
			return p, nil
		}
	}
	return nil, errors.New("not authorized")
}

func (m *MockProjectRepo) Update(userName string, id uint, project *models.Project) (*models.Project, error) {
	existing, err := m.Find(userName, id)
	if err != nil {
		return nil, err
	}
	existing.Title = project.Title
	existing.Description = project.Description
	existing.Platform = project.Platform
	existing.Client = project.Client
	existing.EstimatedFee = project.EstimatedFee
	existing.Status = project.Status
	existing.Deadline = project.Deadline
	existing.UpdatedAt = time.Now()
	return existing, nil
}

func (m *MockProjectRepo) SoftDelete(userName string, id uint) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 内部呼び出しでデッドロックを防ぐ
	var found *models.Project
	for _, p := range m.userIndex[userName] {
		if p.ID == id {
			found = p
			break
		}
	}
	if found == nil {
		return errors.New("not authorized or not found")
	}
	found.DeletedAt.Valid = true
	found.DeletedAt.Time = time.Now()
	return nil

}

func (m *MockProjectRepo) HardDelete(id uint) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.projects, id)
	return nil
}

func (m *MockProjectRepo) FindAllForOwner() ([]models.Project, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var list []models.Project
	for _, p := range m.projects {
		list = append(list, *p)
	}
	return list, nil
}
