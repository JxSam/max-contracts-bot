package kb

import (
	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

type Buttons struct {
	Number string
	Active bool
}

func boolToText(v bool) (string, string) {
	if v {
		return "Отключить", "выкл"
	}
	return "Включить", "вкл"
}

func MainMenuButton(Builder *maxbot.Keyboard) *maxbot.Keyboard {
	k := Builder
	k.AddRow().AddCallback("Главное меню", schemes.NEGATIVE, "/menu")

	return k
}

func MainKeyboard(Builder *maxbot.Keyboard) *maxbot.Keyboard {
	k := Builder
	k.
		AddRow().
		AddCallback("📋 Список контрактов", schemes.POSITIVE, "/notifier").
		AddCallback("❔ Помощь", schemes.POSITIVE, "/help")
	k.AddRow().AddCallback("➕ Добавить контракт", schemes.POSITIVE, "/contract")

	return k
}

func NotifierKeyboard(Builder *maxbot.Keyboard, buttons []Buttons) *maxbot.Keyboard {
	k := Builder
	for _, b := range buttons {
		state, cmd := boolToText(b.Active)
		k.
			AddRow().
			AddLink("№"+b.Number, schemes.POSITIVE, "https://zakupki.gov.ru/epz/contract/contractCard/common-info.html?reestrNumber="+b.Number).
			AddCallback(state, schemes.POSITIVE, cmd+" "+b.Number)
	}
	k.AddRow().AddCallback("Главное меню", schemes.NEGATIVE, "/menu")
	return k
}
