package repository

import (
	"context"
	"fmt"

	"github.com/JxSam/max-contracts-bot/internal/bot/storage"
	"github.com/JxSam/max-contracts-bot/internal/model"
	"github.com/JxSam/max-contracts-bot/pkg/database"
	"github.com/JxSam/max-contracts-bot/pkg/sqlstore"
)

type UserRepository struct {
	pgContext database.Interface
	sqlStore  sqlstore.Interface
}

func NewUserRepository(pgContext database.Interface, sqlStore sqlstore.Interface) *UserRepository {
	return &UserRepository{
		pgContext: pgContext,
		sqlStore:  sqlStore,
	}
}

func (r *UserRepository) CreateUser(tx context.Context, chatID int64) error {
	query, err := r.sqlStore.GetQuery("create_user.sql")
	if err != nil {
		return err
	}

	ctx, cancel := r.pgContext.DefaultTimeoutCtx()
	defer cancel()

	_, err = r.pgContext.TxOrDb(tx).Exec(ctx, query, chatID)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetUserByID(tx context.Context, chatID int64) (*model.User, error) {
	query, err := r.sqlStore.GetQuery("get_user_by_id.sql")
	if err != nil {
		return nil, err
	}

	ctx, cancel := r.pgContext.DefaultTimeoutCtx()
	defer cancel()

	var user model.User

	if err := r.pgContext.TxOrDb(tx).QueryRow(ctx, query, chatID).Scan(
		&user.ChatID,
		&user.Name,
		&user.Verification,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByUsername(username string) *model.User {

	user := model.User{
		Name: &username,
	}

	return &user
}

func (r *UserRepository) UpdateUserByID(tx context.Context, user *model.User) (*model.User, error) {
	query, err := r.sqlStore.GetQuery("update_user.sql")
	if err != nil {
		return nil, err
	}

	ctx, cancel := r.pgContext.DefaultTimeoutCtx()
	defer cancel()

	_, err = r.pgContext.TxOrDb(tx).Exec(ctx, query, &user.ChatID,
		&user.Name,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetAllUsers() ([]storage.User, error) {
	query, err := r.sqlStore.GetQuery("get_users.sql")
	if err != nil {
		return nil, fmt.Errorf("SQL query get_users not found")
	}

	var users []storage.User

	ctx, cancel := r.pgContext.DefaultTimeoutCtx()
	defer cancel()

	rows, err := r.pgContext.TxOrDb(ctx).Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	return users, nil
}

func (r *UserRepository) GetUsers() ([]model.UserDefault, error) {
	query, err := r.sqlStore.GetQuery("get_all_users.sql")
	if err != nil {
		return nil, fmt.Errorf("SQL query get_contracts not found")
	}

	var users []model.UserDefault

	ctx, cancel := r.pgContext.DefaultTimeoutCtx()
	defer cancel()

	rows, err := r.pgContext.TxOrDb(ctx).Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user model.UserDefault
		err := rows.Scan(&user.ChatID)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при чтении строк: %w", err)
	}

	return users, nil
}
