package transaction

import (
	"context"
	"time"

	"github.com/JxSam/max-contracts-bot/pkg/database"
)

type UserService interface {
	CreateUser(ctx context.Context, chatID int64) error
}

type UsersContractsService interface {
	CreateUsersContracts(chatID int64, contractID int64) error
}

type ContractsService interface {
	CreateContract(id int64, link string, created_at time.Time, description string, name_product string, defaults bool, link_contract string) error
}

type Transaction struct {
	userService           UserService
	contractsService      ContractsService
	usersContractsService UsersContractsService
	txManager             database.TransactionManager
}

func NewTransaction(us UserService, cs ContractsService, ucs UsersContractsService, tx database.TransactionManager) *Transaction {
	return &Transaction{
		userService:           us,
		contractsService:      cs,
		usersContractsService: ucs,
		txManager:             tx,
	}
}
