package repository

import (
	"errors"
)

var (
	Admin *Manager
)

func init() {
	Admin = NewManager()
}

type Manager struct {
	repositories map[string]interface{}
}

func NewManager() *Manager {
	return &Manager{
		repositories: make(map[string]interface{}),
	}
}

func L(name string) interface{} {
	return Admin.LoadRepository(name)
}

func (m *Manager) LoadRepository(name string) interface{} {
	return m.repositories[name]
}

func (m *Manager) Register(name string, repository interface{}) error {
	if _, ok := m.repositories[name]; ok {
		return errors.New("already exist repository: " + name)
	}
	m.repositories[name] = repository
	return nil
}
