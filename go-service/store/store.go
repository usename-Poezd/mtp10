package store

import (
	"sync"

	"go-service/models"

	"github.com/google/uuid"
)

type Store struct {
	mu    sync.RWMutex
	todos map[string]models.Todo
}

func NewStore() *Store {
	return &Store{todos: make(map[string]models.Todo)}
}

func (s *Store) GetAll() []models.Todo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]models.Todo, 0, len(s.todos))
	for _, todo := range s.todos {
		result = append(result, todo)
	}

	if result == nil {
		return []models.Todo{}
	}

	return result
}

func (s *Store) GetByID(id string) (models.Todo, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	todo, ok := s.todos[id]
	return todo, ok
}

func (s *Store) Create(title string) models.Todo {
	s.mu.Lock()
	defer s.mu.Unlock()

	todo := models.Todo{
		ID:        uuid.New().String(),
		Title:     title,
		Completed: false,
	}
	s.todos[todo.ID] = todo

	return todo
}

func (s *Store) Update(id, title string, completed bool) (models.Todo, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	todo, ok := s.todos[id]
	if !ok {
		return models.Todo{}, false
	}

	todo.Title = title
	todo.Completed = completed
	s.todos[id] = todo

	return todo, true
}

func (s *Store) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.todos[id]; !ok {
		return false
	}

	delete(s.todos, id)
	return true
}
