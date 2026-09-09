package handler

import (
	"errors"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type userInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type userStore struct {
	sync.RWMutex
	users  map[string]User
	nextID int
}

// ponytail: process-local mock; use a database when data must survive cold starts or multiple Lambda instances.
var mockUsers = newMockUserStore()

func newMockUserStore() *userStore {
	return &userStore{
		users: map[string]User{
			"1": {ID: "1", Name: "Alice", Email: "alice@example.com"},
			"2": {ID: "2", Name: "Bob", Email: "bob@example.com"},
		},
		nextID: 3,
	}
}

func bindUser(c *gin.Context) (userInput, error) {
	var input userInput
	if err := c.ShouldBindJSON(&input); err != nil {
		return userInput{}, err
	}
	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.TrimSpace(input.Email)
	if input.Name == "" || input.Email == "" {
		return userInput{}, errors.New("name and email are required")
	}
	return input, nil
}

func (s *userStore) list() []User {
	s.RLock()
	defer s.RUnlock()

	ids := make([]string, 0, len(s.users))
	for id := range s.users {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	users := make([]User, 0, len(ids))
	for _, id := range ids {
		users = append(users, s.users[id])
	}
	return users
}

func (s *userStore) get(id string) (User, bool) {
	s.RLock()
	defer s.RUnlock()
	user, ok := s.users[id]
	return user, ok
}

func (s *userStore) create(input userInput) User {
	s.Lock()
	defer s.Unlock()

	id := strconv.Itoa(s.nextID)
	user := User{ID: id, Name: input.Name, Email: input.Email}
	s.users[id] = user
	s.nextID++
	return user
}

func (s *userStore) update(id string, input userInput) (User, bool) {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.users[id]; !ok {
		return User{}, false
	}
	user := User{ID: id, Name: input.Name, Email: input.Email}
	s.users[id] = user
	return user, true
}

func (s *userStore) delete(id string) bool {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.users[id]; !ok {
		return false
	}
	delete(s.users, id)
	return true
}

func writeError(c *gin.Context, status int, message string) {
	c.JSON(status, map[string]string{"error": message})
}
