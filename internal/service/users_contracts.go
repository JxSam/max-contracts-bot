package service

import (
	"log/slog"
	"strings"

	"github.com/JxSam/max-contracts-bot/internal/bot/storage"
)

type UsersContractsRepository interface {
	CreateUsersContracts(chatID int64, contractID int64) error
	UpdateUsersContracts(userID int64, col string, flag bool) error
	CheckUsersContracts(userID int64) (*storage.UsersContracts, error)
	GetUsersContracts(userID int64) ([]storage.UsersContracts, error)
	GetUsersActiveContracts(userID int64) ([]storage.UsersContracts, error)
	UpdateActiveUsersContracts(userID int64, flag bool, message string) error
}

type UsersContractsService struct {
	usersContractsRepository UsersContractsRepository
	messages                 map[string]string
	log                      *slog.Logger
}

func NewUsersContractsService(UsersContractsRepository UsersContractsRepository, log *slog.Logger) *UsersContractsService {
	return &UsersContractsService{
		usersContractsRepository: UsersContractsRepository,
		log:                      log,
	}
}

func (s *UsersContractsService) CreateUsersContracts(chatID int64, contractID int64) error {
	return s.usersContractsRepository.CreateUsersContracts(chatID, contractID)
}

func (s *UsersContractsService) CheckUsersContracts(userID int64) (*storage.UsersContracts, error) {
	return s.usersContractsRepository.CheckUsersContracts(userID)
}

func (s *UsersContractsService) GetUsersContracts(userID int64) ([]storage.UsersContracts, error) {
	return s.usersContractsRepository.GetUsersContracts(userID)
}

func (s *UsersContractsService) GetUsersActiveContractsService(userID int64) ([]storage.UsersContracts, error) {
	return s.usersContractsRepository.GetUsersActiveContracts(userID)
}

func (s *UsersContractsService) UpdateActiveUsersContracts(userID int64, message string) error {
	sp := strings.Split(message, " ")
	var flag bool

	switch sp[0] {
	case "вкл":
		flag = true
	case "выкл":
		flag = false
	default:
		return nil
	}
	return s.usersContractsRepository.UpdateActiveUsersContracts(userID, flag, message)
}
