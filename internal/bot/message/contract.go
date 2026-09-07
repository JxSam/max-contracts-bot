package executor

import (
	"fmt"

	"github.com/JxSam/max-contracts-bot/internal/repository"
)

const message_contract = `
Контракт № %v обновился!

Данные:
- Дата события: %v
- Описание: %v
- Предмет контракта: %v`

func Contract(d *repository.Contracts) *string {
	time := d.CreatedAt.Format("2006-01-02 15:04:05")

	newMessage := fmt.Sprintf(message_contract,
		d.Id, time, d.Description, d.NameProduct)

	return &newMessage
}
