package notifier

import (
	"context"
	"log/slog"
	"time"

	"github.com/JxSam/max-contracts-bot/internal/bot/storage"
	"github.com/JxSam/max-contracts-bot/internal/model"
	"github.com/JxSam/max-contracts-bot/internal/repository"
)

type UserService interface {
	GetAllUsers()
	GetUserByID(chatID int64) (*model.User, error)
	GetUsers() ([]model.UserDefault, error)
}

type Notifier struct {
	log              *slog.Logger
	ch               chan model.Notify
	userService      UserService
	contractsService ContractsService
	active_users     UsersContractsService
}

type ContractsService interface {
	GetAllContractsService() ([]repository.Contracts, error)
	UpdateDefaultContract(ctx context.Context, id int64, defaults bool) error
	CreateContract(id int64, link string, created_at time.Time, description string, name_product string, defaults bool, link_contract string) error
	UpdateContractService(
		ctx context.Context,
		id int64,
		date time.Time,
		description string,
		name_product string,
	) error
}

type UsersContractsService interface {
	GetUsersActiveContractsService(UserID int64) ([]storage.UsersContracts, error)
}

func New(
	log *slog.Logger,
	ch chan model.Notify,
	userService UserService,
	contractService ContractsService,
	active_users UsersContractsService,
) *Notifier {
	return &Notifier{
		log:              log,
		ch:               ch,
		userService:      userService,
		contractsService: contractService,
		active_users:     active_users,
	}
}

func (n *Notifier) Run(ctx context.Context) {
	go n.updateContracts(ctx)
	go n.updateListContracts(ctx)
}
