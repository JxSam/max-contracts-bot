package bot

import (
	"fmt"
	"strconv"
	"time"

	"github.com/JxSam/max-contracts-bot/internal/parser"

	"github.com/JxSam/max-contracts-bot/internal/bot/command"
	kb "github.com/JxSam/max-contracts-bot/internal/bot/keyboard"
	"github.com/JxSam/max-contracts-bot/internal/bot/storage"
	"github.com/JxSam/max-contracts-bot/internal/model"
)

type Transaction interface {
	CreateUser(chatID int64) error
}

type UserService interface {
	GetUserByID(chatID int64) (*model.User, error)
	UpdateUserByID(user *model.User) (*model.User, error)
	GetUserByUsername(username string) *model.User
}
type ContractsService interface {
	CreateContract(id int64, link string, created_at time.Time, description string, name_product string, defaults bool, link_contract string) error
}

type UsersContractsService interface {
	CreateUsersContracts(chatID int64, contractID int64) error
	UpdateActiveUsersContracts(userID int64, message string) error
	GetUsersActiveContractsService(UserID int64) ([]storage.UsersContracts, error)
	GetUsersContracts(userID int64) ([]storage.UsersContracts, error)
}

type Handler struct {
	transaction           Transaction
	userService           UserService
	contractsService      ContractsService
	usersContractsService UsersContractsService
}

func NewHandler(transaction Transaction, userService UserService, contractsService ContractsService, usersContractsService UsersContractsService) *Handler {
	return &Handler{
		transaction:           transaction,
		userService:           userService,
		contractsService:      contractsService,
		usersContractsService: usersContractsService,
	}
}

func (h *Handler) StartHandler(chatID int64) string {
	err := h.transaction.CreateUser(chatID)
	if err != nil {
		fmt.Println(err)
	}

	return command.Start
}

func (h *Handler) NotifierHandler(chatID int64) (string, []kb.Buttons) {
	var us []storage.UsersContracts
	var kbs []kb.Buttons
	us, _ = h.usersContractsService.GetUsersContracts(chatID)
	for _, user := range us {
		kbs = append(kbs, kb.Buttons{
			Number: strconv.FormatInt(user.ContractID, 10),
			Active: user.Active,
		})
	}
	return command.Notifier, kbs
}

func (h *Handler) CreateContractHandler(chatID int64, id int64, link string) string {
	var Event parser.Event
	var err error

	if id <= 0 {
		return "Ошибка: некорректный ID контракта."
	}

	Event, err = parser.ParseLink(strconv.FormatInt(id, 10), false)
	if err != nil {
		return err.Error()
	}
	err = h.contractsService.CreateContract(id, link, Event.EventTime, Event.Description, Event.ItemName, false, link)

	err = h.usersContractsService.CreateUsersContracts(chatID, id)

	if err != nil {
		fmt.Println(err)
	}

	return "Контракт: " + strconv.FormatInt(id, 10) + " успешно добавлен!"
}

func (h *Handler) NewContractHandler(chatID int64) string {
	return command.ContractLink
}

func (h *Handler) MainMenuHandler(chatID int64) string {
	return command.Start
}

func (h *Handler) HelpHandler(chatID int64) string {
	return command.Help
}

func (h *Handler) UpdateActiveUsersContracts(userID int64, message string) string {
	err := h.usersContractsService.UpdateActiveUsersContracts(userID, message)
	fmt.Println(err)
	return "Уведомление о контракте отключено"
}
