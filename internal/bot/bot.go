package bot

import (
	"context"
	"fmt"
	"log/slog"

	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"

	"github.com/JxSam/max-contracts-bot/internal/config"
	"github.com/JxSam/max-contracts-bot/internal/model"
)

type Bot struct {
	api    *maxbot.Api
	router *Router
	cfg    config.BotConfig
	log    *slog.Logger
}

func NewBot(cfg config.BotConfig, log *slog.Logger, router *Router) (*Bot, error) {
	api, err := maxbot.New(cfg.Token)
	if err != nil {
		return nil, err
	}
	info, err := api.Bots.GetBot(context.Background())
	fmt.Printf("Get me: %#v %#v", info, err)

	return &Bot{
		api:    api,
		router: router,
		cfg:    cfg,
		log:    log,
	}, nil
}

func (b *Bot) Run(ctx context.Context) {
	for upd := range b.api.GetUpdates(ctx) {
		switch upd := upd.(type) {
		case *schemes.MessageCreatedUpdate:
			message, keyboard := b.router.Routing(&upd.Message, b)
			if keyboard == nil {
				err := b.api.Messages.Send(ctx, maxbot.NewMessage().SetChat(upd.Message.Recipient.ChatId).SetText(message))
				if err != nil {
					//b.log.Error("error", slog.Any("error on message sending: ", err))
				}
			} else {
				_ = b.api.Messages.Send(ctx, maxbot.NewMessage().SetChat(upd.Message.Recipient.ChatId).AddKeyboard(keyboard).SetText(message))
			}
		case *schemes.BotStartedUpdate:
			v := schemes.Message{
				Recipient: schemes.Recipient{
					ChatId: upd.ChatId,
				},
				Body: schemes.MessageBody{
					Text: "/start",
				},
			}
			message, keyboard := b.router.Routing(&v, b)
			if keyboard == nil {
				err := b.api.Messages.Send(ctx, maxbot.NewMessage().SetChat(v.Recipient.ChatId).SetText(message))
				if err != nil {
					b.log.Error("error", slog.Any("error on message sending: ", err))
				}
			} else {
				_ = b.api.Messages.Send(ctx, maxbot.NewMessage().SetChat(v.Recipient.ChatId).AddKeyboard(keyboard).SetText(message))
			}
		case *schemes.MessageCallbackUpdate:
			v := schemes.Message{
				Recipient: schemes.Recipient{
					ChatId: upd.Message.Recipient.ChatId,
				},
				Body: schemes.MessageBody{
					Text: upd.Callback.Payload,
				},
			}
			message, keyboard := b.router.Routing(&v, b)
			msg := maxbot.NewMessage().
				SetChat(upd.Message.Recipient.ChatId).
				SetText(message)
			if keyboard != nil {
				msg.AddKeyboard(keyboard)
			}
			if upd.Callback.Payload == "/start" {
				err := b.api.Messages.Send(ctx, msg)
				if err != nil {
					b.log.Error("error", slog.Any("error on message sending", err))
				}
			} else {
				err := b.api.Messages.EditMessage(ctx, upd.Message.Body.Mid, msg)
				if err != nil {
					b.log.Error("error", slog.Any("error on message editing", err))
				}
			}
		}
	}
}

func (b *Bot) Notify(ctx context.Context, ch chan model.Notify) {
	var err error
	for {
		select {
		case <-ctx.Done():
			return
		case notify, ok := <-ch:
			if !ok {
				return
			}

			if notify.LinkButton != nil {
				keyboard := b.api.Messages.NewKeyboardBuilder()
				keyboard.
					AddRow().
					AddLink(notify.LinkButton.Text, schemes.POSITIVE, notify.LinkButton.Link)
				err = b.api.Messages.Send(ctx, maxbot.NewMessage().SetChat(notify.UserID).AddKeyboard(keyboard).SetText(notify.Message))
			} else {
				err = b.api.Messages.Send(ctx, maxbot.NewMessage().SetChat(notify.UserID).SetText(notify.Message))
			}

			if err != nil {
				b.log.Error("error", slog.Any("error on message sending: ", err))
			}
		}
	}
}
