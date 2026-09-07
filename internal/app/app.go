package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/JxSam/max-contracts-bot/internal/bot"
	"github.com/JxSam/max-contracts-bot/internal/bot/notifier"
	"github.com/JxSam/max-contracts-bot/internal/config"
	"github.com/JxSam/max-contracts-bot/internal/model"
	"github.com/JxSam/max-contracts-bot/internal/repository"
	"github.com/JxSam/max-contracts-bot/internal/service"
	"github.com/JxSam/max-contracts-bot/internal/transaction"
	"github.com/JxSam/max-contracts-bot/pkg/database"
	"github.com/JxSam/max-contracts-bot/pkg/logger"
	"github.com/JxSam/max-contracts-bot/pkg/sqlstore"
)

func Run(cfg *config.Config) error {
	log := logger.New(cfg.Env)

	pgContext, err := database.New(&cfg.DatabaseConfig)
	if err != nil {
		return err
	}

	sqlstore, err := sqlstore.New(cfg.SqlConfig.Path)
	if err != nil {
		return err
	}

	// Repositories
	userRepository := repository.NewUserRepository(pgContext, sqlstore)
	usersContractsRepository := repository.NewUsersContractsRepository(pgContext, sqlstore)
	contractsRepository := repository.NewContractsRepository(pgContext, sqlstore)

	// Services
	userService := service.NewUserService(userRepository)
	usersContractsService := service.NewUsersContractsService(usersContractsRepository, log)
	contractsService := service.NewContractsService(contractsRepository, log)

	// Transactions
	transaction := transaction.NewTransaction(userService, contractsService, usersContractsService, pgContext)

	notify := make(chan model.Notify)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Bot
	handler := bot.NewHandler(transaction, userService, contractsService, usersContractsService)
	router := bot.NewRouter(handler)
	bot, err := bot.NewBot(cfg.BotConfig, log, router)
	if err != nil {
		return err
	}

	go bot.Run(ctx)
	go bot.Notify(ctx, notify)

	// Notifier

	notifier := notifier.New(log, notify, userService, contractsService, usersContractsService)
	go notifier.Run(ctx)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	return nil
}
