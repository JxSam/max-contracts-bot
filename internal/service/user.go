package service

import (
	"context"

	"github.com/JxSam/max-contracts-bot/internal/bot/storage"
	"github.com/JxSam/max-contracts-bot/internal/model"
)

type UserRepository interface {
	CreateUser(tx context.Context, chatID int64) error
	UpdateUserByID(tx context.Context, user *model.User) (*model.User, error)
	GetUserByID(tx context.Context, chatID int64) (*model.User, error)
	GetAllUsers() ([]storage.User, error)
	GetUserByUsername(username string) *model.User
	GetUsers() ([]model.UserDefault, error)
}

type UserService struct {
	userRepository UserRepository
}

func NewUserService(userRepository UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) CreateUser(tx context.Context, chatID int64) error {
	return s.userRepository.CreateUser(tx, chatID)
}

func (s *UserService) UpdateUserByID(user *model.User) (*model.User, error) {
	return s.userRepository.UpdateUserByID(context.Background(), user)
}

func (s *UserService) GetUserByID(chatID int64) (*model.User, error) {
	return s.userRepository.GetUserByID(context.Background(), chatID)
}

func (s *UserService) GetUserByUsername(username string) *model.User {
	return s.userRepository.GetUserByUsername(username)
}

func (s *UserService) GetUsers() ([]model.UserDefault, error) {
	return s.userRepository.GetUsers()
}

func (us *UserService) GetAllUsers() {
	users, err := us.userRepository.GetAllUsers()
	if err != nil {
		return
	}

	storage.Update(users)
}
