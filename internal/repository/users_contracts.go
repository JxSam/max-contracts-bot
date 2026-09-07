package repository

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/JxSam/max-contracts-bot/internal/bot/storage"
	"github.com/JxSam/max-contracts-bot/pkg/database"
	"github.com/JxSam/max-contracts-bot/pkg/sqlstore"
)

type UsersContractsRepository struct {
	pgContext database.Interface
	sqlStore  sqlstore.Interface
}

type UserContract struct {
	UserID     int64
	ContractID int64
	Active     bool
}

func NewUsersContractsRepository(pgContext database.Interface, sqlStore sqlstore.Interface) *UsersContractsRepository {
	return &UsersContractsRepository{
		pgContext: pgContext,
		sqlStore:  sqlStore,
	}
}

func (r *UsersContractsRepository) CreateUsersContracts(chatID int64, contractID int64) error {
	query, err := r.sqlStore.GetQuery("create_users_contracts.sql")
	if err != nil {
		return err
	}

	ctx, cancel := r.pgContext.DefaultTimeoutCtx()
	defer cancel()

	_, err = r.pgContext.TxOrDb(ctx).Exec(ctx, query, chatID, contractID)
	if err != nil {
		return err
	}

	return nil
}

func (r *UsersContractsRepository) UpdateUsersContracts(userID int64, col string, flag bool) error {
	ctx, cancel := r.pgContext.DefaultTimeoutCtx()
	defer cancel()

	query := fmt.Sprintf("UPDATE users_contracts SET %v = %v WHERE user_id = '%d'", col, flag, userID)

	_, err := r.pgContext.TxOrDb(ctx).Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("error update notify: %w", err)
	}

	return nil
}

func (r *UsersContractsRepository) CheckUsersContracts(userID int64) (*storage.UsersContracts, error) {
	query, err := r.sqlStore.GetQuery("get_users_contracts.sql")
	if err != nil {
		return nil, fmt.Errorf("SQL query get_users_contracts not found")
	}

	ctx, cancel := r.pgContext.DefaultTimeoutCtx()
	defer cancel()

	var users_contracts storage.UsersContracts

	err = r.pgContext.TxOrDb(ctx).QueryRow(ctx, query, userID).Scan(
		&users_contracts.ChatID,
		&users_contracts.ContractID,
		&users_contracts.Active,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	return &users_contracts, nil
}

func (r *UsersContractsRepository) UpdateActiveUsersContracts(userID int64, flag bool, message string) error {
	re := regexp.MustCompile(`\d+`)
	contractStr := re.FindString(message)
	if contractStr == "" {
		return fmt.Errorf("contract id not found")
	}
	contractID, err := strconv.ParseInt(contractStr, 10, 64)
	if err != nil {
		return err
	}
	fmt.Println(contractID)
	ctx, cancel := r.pgContext.DefaultTimeoutCtx()
	defer cancel()

	query := "UPDATE users_contracts SET active = $1 WHERE user_id = $2 AND contract_id = $3"

	_, err = r.pgContext.TxOrDb(ctx).Exec(ctx, query, flag, userID, contractID)
	if err != nil {
		return fmt.Errorf("error update notify: %w", err)
	}

	return nil
}

func (r *UsersContractsRepository) GetUsersActiveContracts(
	contractID int64,
) ([]storage.UsersContracts, error) {

	query, err := r.sqlStore.GetQuery("get_users_active_contracts.sql")
	if err != nil {
		return nil, fmt.Errorf("SQL query get_users not found")
	}

	var activeUsers []storage.UsersContracts

	ctx, cancel := r.pgContext.DefaultTimeoutCtx()
	defer cancel()

	rows, err := r.pgContext.TxOrDb(ctx).Query(
		ctx,
		query,
		contractID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {

		var user storage.UsersContracts

		err := rows.Scan(
			&user.ChatID,
			&user.ContractID,
			&user.Active,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to scan row: %w",
				err,
			)
		}

		activeUsers = append(activeUsers, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"rows iteration error: %w",
			err,
		)
	}

	return activeUsers, nil
}

func (r *UsersContractsRepository) GetUsersContracts(
	chatID int64,
) ([]storage.UsersContracts, error) {

	query, err := r.sqlStore.GetQuery("get_users_contracts_2.sql")
	if err != nil {
		return nil, fmt.Errorf("SQL query get_users not found")
	}

	var activeUsers []storage.UsersContracts

	ctx, cancel := r.pgContext.DefaultTimeoutCtx()
	defer cancel()

	rows, err := r.pgContext.TxOrDb(ctx).Query(
		ctx,
		query,
		chatID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {

		var user storage.UsersContracts

		err := rows.Scan(
			&user.ChatID,
			&user.ContractID,
			&user.Active,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to scan row: %w",
				err,
			)
		}

		activeUsers = append(activeUsers, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"rows iteration error: %w",
			err,
		)
	}

	return activeUsers, nil
}
