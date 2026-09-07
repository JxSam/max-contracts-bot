package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/JxSam/max-contracts-bot/internal/repository"
)

type ContractsRepository interface {
	CreateContract(id int64, link string, created_at time.Time, description string, name_product string, defaults bool, link_contract string) error
	GetAllContracts() ([]repository.Contracts, error)
	UpdateContract(ctx context.Context, id int64, date time.Time, description string, name_product string) error
	UpdateDefaultContract(ctx context.Context, id int64, defaults bool) error
}

type ContractsService struct {
	contractsRepository ContractsRepository
	log                 *slog.Logger
}

func NewContractsService(contractsRepository ContractsRepository, log *slog.Logger) *ContractsService {
	return &ContractsService{
		contractsRepository: contractsRepository,
		log:                 log,
	}
}

func (s *ContractsService) CreateContract(id int64, link string, created_at time.Time, description string, name_product string, defaults bool, link_contract string) error {
	return s.contractsRepository.CreateContract(id, link, created_at, description, name_product, defaults, link_contract)
}

func (s *ContractsService) GetAllContractsService() ([]repository.Contracts, error) {
	return s.contractsRepository.GetAllContracts()
}

func (s *ContractsService) UpdateContractService(ctx context.Context, id int64, date time.Time, description string, name_product string) error {
	return s.contractsRepository.UpdateContract(ctx, id, date, description, name_product)
}

func (s *ContractsService) UpdateDefaultContract(ctx context.Context, id int64, defaults bool) error {
	return s.contractsRepository.UpdateDefaultContract(ctx, id, defaults)
}
