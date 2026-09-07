package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/JxSam/max-contracts-bot/pkg/database"
	"github.com/JxSam/max-contracts-bot/pkg/sqlstore"
)

type ContractsRepository struct {
	pgContext database.Interface
	sqlStore  sqlstore.Interface
}

type Contracts struct {
	Id           int64     `json:"user_id"`
	Link         string    `json:"link"`
	CreatedAt    time.Time `json:"created_at"`
	Description  string    `json:"description"`
	NameProduct  string    `json:"nameProduct"`
	Default      bool      `json:"default"`
	LinkContract string    `json:"link_contract"`
}

func NewContractsRepository(pgContext database.Interface, sqlStore sqlstore.Interface) *ContractsRepository {
	return &ContractsRepository{
		pgContext: pgContext,
		sqlStore:  sqlStore,
	}
}

func (r *ContractsRepository) CreateContract(id int64, link string, created_at time.Time, description string, name_product string, defaults bool, link_contract string) error {
	query, err := r.sqlStore.GetQuery("create_contracts.sql")
	if err != nil {
		return err
	}

	ctx, cancel := r.pgContext.DefaultTimeoutCtx()
	defer cancel()

	_, err = r.pgContext.TxOrDb(ctx).Exec(ctx, query, id, link, created_at, description, name_product, defaults, link_contract)
	if err != nil {
		return err
	}

	return nil
}

func (r *ContractsRepository) UpdateContract(
	ctx context.Context,
	id int64,
	date time.Time,
	description string,
	name_product string,
) error {
	query := "UPDATE contracts SET created_at = $1, description = $3, name_product = $4 WHERE id = $2"

	ctx, cancel := r.pgContext.DefaultTimeoutCtx()
	defer cancel()

	_, err := r.pgContext.TxOrDb(ctx).Exec(ctx, query, date, id, description, name_product)
	if err != nil {
		return fmt.Errorf("error update contract date: %w", err)
	}

	return nil
}

func (r *ContractsRepository) UpdateDefaultContract(
	ctx context.Context,
	id int64,
	defaults bool,
) error {
	query, err := r.sqlStore.GetQuery("update_contract.sql")
	if err != nil {
		return fmt.Errorf("SQL query update_contract.sql not found")
	}
	ctx, cancel := r.pgContext.DefaultTimeoutCtx()
	defer cancel()

	_, err = r.pgContext.TxOrDb(ctx).Exec(ctx, query, defaults, id)
	if err != nil {
		return fmt.Errorf("error update contract date: %w", err)
	}

	return nil
}

func (r *ContractsRepository) GetAllContracts() ([]Contracts, error) {
	query, err := r.sqlStore.GetQuery("get_contracts.sql")
	if err != nil {
		return nil, fmt.Errorf("SQL query get_contracts.sql not found")
	}

	var contracts []Contracts

	ctx, cancel := r.pgContext.DefaultTimeoutCtx()
	defer cancel()

	rows, err := r.pgContext.TxOrDb(ctx).Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var contract Contracts
		err := rows.Scan(&contract.Id, &contract.Link, &contract.CreatedAt, &contract.Description, &contract.NameProduct, &contract.Default, &contract.LinkContract)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		contracts = append(contracts, contract)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при чтении строк: %w", err)
	}

	return contracts, nil
}
