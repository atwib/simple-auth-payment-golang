package repository

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"sync"

	"test-mnc/model"
)

type AuthRepository struct {
	mu     sync.Mutex
	dbFile string
	users  []model.User
}

func NewAuthRepository(dbFile string) *AuthRepository {
	repo := &AuthRepository{dbFile: dbFile}
	repo.loadUsers()
	return repo
}

func (r *AuthRepository) loadUsers() error {
	if _, err := os.Stat(r.dbFile); os.IsNotExist(err) {
		if err := ioutil.WriteFile(r.dbFile, []byte("[]"), 0644); err != nil {
			return err
		}
	}

	data, err := ioutil.ReadFile(r.dbFile)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &r.users)
}

func (r *AuthRepository) saveUsers() error {
	data, err := json.MarshalIndent(r.users, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(r.dbFile, data, 0644)
}

func (r *AuthRepository) AddUser(user model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.userExists(user.Username) {
		return errors.New("username already exists")
	}

	r.users = append(r.users, user)
	return r.saveUsers()
}

func (r *AuthRepository) userExists(username string) bool {
	for _, u := range r.users {
		if u.Username == username {
			return true
		}
	}
	return false
}

func (r *AuthRepository) GetUserByUsername(username string) (*model.User, error) {
	for _, u := range r.users {
		if u.Username == username {
			return &u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *AuthRepository) GetByUserId(userId string) (*model.User, error) {
	for _, u := range r.users {
		if u.UserId == userId {
			return &u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *AuthRepository) UpdateUser(user *model.User) error {
	for i, u := range r.users {
		if u.Username == user.Username {
			r.users[i] = *user
			return r.saveUsers()
		}
	}
	return errors.New("user not found")
}
