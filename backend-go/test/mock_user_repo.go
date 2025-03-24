package test

import (
	"errors"
	"sync"
	"time"

	"github.com/nakachan-ing/fukuworks/backend-go/internal/domain/models"
	"github.com/nakachan-ing/fukuworks/backend-go/internal/domain/repositories"
)

type MockUserRepo struct {
	mu    sync.Mutex
	users map[string]*models.User // key = user name
	idSeq uint
}

func NewMockUserRepo() repositories.UserRepository {
	return &MockUserRepo{
		users: make(map[string]*models.User),
		idSeq: 1,
	}
}

func (m *MockUserRepo) Create(user *models.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.users[user.Name]; exists {
		return errors.New("user already exists")
	}
	user.ID = m.idSeq
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	m.idSeq++
	m.users[user.Name] = user
	return nil
}

func (m *MockUserRepo) Find(name string) (*models.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	user, ok := m.users[name]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *MockUserRepo) Update(name string, updated *models.User) (*models.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	existing, ok := m.users[name]
	if !ok {
		return nil, errors.New("user not found")
	}
	existing.Name = updated.Name
	existing.Email = updated.Email
	existing.Password = updated.Password
	existing.UpdatedAt = time.Now()
	return existing, nil
}

func (m *MockUserRepo) SoftDelete(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	user, ok := m.users[name]
	if !ok {
		return errors.New("user not found")
	}
	user.DeletedAt.Valid = true
	user.DeletedAt.Time = time.Now()
	return nil
}

func (m *MockUserRepo) FindAll() ([]models.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var result []models.User
	for _, u := range m.users {
		result = append(result, *u)
	}
	return result, nil
}

func (m *MockUserRepo) HardDelete(id uint) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for name, u := range m.users {
		if u.ID == id {
			delete(m.users, name)
			return nil
		}
	}
	return errors.New("user not found")
}
