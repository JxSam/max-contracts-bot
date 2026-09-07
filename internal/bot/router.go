package bot

import (
	"strconv"
	"strings"

	kb "github.com/JxSam/max-contracts-bot/internal/bot/keyboard"
	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

type Router struct {
	handler *Handler
}

func NewRouter(handler *Handler) *Router {
	return &Router{
		handler: handler,
	}
}

func (r *Router) Routing(message *schemes.Message, b *Bot) (string, *maxbot.Keyboard) {
	switch message.Body.Text {
	case "/start":
		message := r.handler.StartHandler(message.Recipient.ChatId)
		k := kb.MainKeyboard(b.api.Messages.NewKeyboardBuilder())
		return message, k
	case "/menu", "меню":
		message := r.handler.MainMenuHandler(message.Recipient.ChatId)
		k := kb.MainKeyboard(b.api.Messages.NewKeyboardBuilder())
		return message, k
	case "/help", "помощь":
		message := r.handler.HelpHandler(message.Recipient.ChatId)
		k := kb.MainMenuButton(b.api.Messages.NewKeyboardBuilder())
		return message, k
	case "/contract", "контракт":
		message := r.handler.NewContractHandler(message.Recipient.ChatId)
		k := kb.MainMenuButton(b.api.Messages.NewKeyboardBuilder())
		return message, k
	case "/list", "уведомления", "/notifier":
		message, buttons := r.handler.NotifierHandler(message.Recipient.ChatId)
		k := kb.NotifierKeyboard(b.api.Messages.NewKeyboardBuilder(), buttons)
		return message, k
	}

	if strings.Contains(message.Body.Text, "zakupki.gov.ru") &&
		strings.Contains(message.Body.Text, "reestrNumber=") {

		parts := strings.Split(message.Body.Text, "reestrNumber=")
		k := kb.MainMenuButton(b.api.Messages.NewKeyboardBuilder())
		if len(parts) > 1 {
			idPart := parts[1]

			if i := strings.Index(idPart, "&"); i != -1 {
				idPart = idPart[:i]
			}

			id, err := strconv.ParseInt(idPart, 10, 64)
			if err != nil {
				return "Некорректный ID контракта!", k
			}

			return r.handler.CreateContractHandler(message.Recipient.ChatId, id, message.Body.Text), kb.MainMenuButton(b.api.Messages.NewKeyboardBuilder())
		}
	}

	if strings.HasPrefix(message.Body.Text, "вкл") || strings.HasPrefix(message.Body.Text, "выкл") {
		return r.handler.UpdateActiveUsersContracts(message.Recipient.ChatId, message.Body.Text), kb.MainMenuButton(b.api.Messages.NewKeyboardBuilder())
	}

	k := kb.MainMenuButton(b.api.Messages.NewKeyboardBuilder())
	return "Неизвестная команда!", k
}
