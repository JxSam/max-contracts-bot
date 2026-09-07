package transaction

import (
	"context"
	"fmt"
)

func (t *Transaction) CreateUsersContracts(chatID int64) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tx, err := t.txManager.BeginTransaction(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer t.txManager.RollbackTransaction(tx)

	// if err := t.usersContractsService.CreateUsersContracts(tx, chatID); err != nil {
	// 	return err
	// }

	if err := t.txManager.CommitTransaction(tx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
