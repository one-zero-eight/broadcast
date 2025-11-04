// Handlers and keyboard UI for the Telegram bot.
// Includes /start command, inline keyboard builder,
// and callback logic for option toggling and selection.
package internal

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// currentOptions keeps in-memory selection state.
// Note: This is shared across all users and chats.
// It is fine for learning, but not for production.
var currentOptions = []bool{false, false, false}

// CallbackHandler handles inline button clicks.
// It toggles options based on callback data, updates the keyboard,
// and on "Select" deletes the original message and sends the result.
func CallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	switch update.CallbackQuery.Data {
	case "btn_opt1":
		currentOptions[0] = !currentOptions[0]
	case "btn_opt2":
		currentOptions[1] = !currentOptions[1]
	case "btn_opt3":
		currentOptions[2] = !currentOptions[2]
	case "btn_select":
		_, err := b.DeleteMessage(ctx, &bot.DeleteMessageParams{
			ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
			MessageID: update.CallbackQuery.Message.Message.ID,
		})
		if err != nil {
			return
		}
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Message.Chat.ID,
			Text:   fmt.Sprintf("Selected options: %v", currentOptions),
		})
		if err != nil {
			return
		}
		return
	}

	b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		MessageID:   update.CallbackQuery.Message.Message.ID,
		ReplyMarkup: EduDegree(),
	})
}

// EduDegree builds an inline keyboard with three toggle buttons
// and one "Select" button. The check mark shows current state.
func EduDegree() models.ReplyMarkup {
	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: buttonText("Option 1", currentOptions[0]), CallbackData: "btn_opt1"},
				{Text: buttonText("Option 2", currentOptions[1]), CallbackData: "btn_opt2"},
				{Text: buttonText("Option 3", currentOptions[2]), CallbackData: "btn_opt3"},
			}, {
				{Text: "Select", CallbackData: "btn_select"},
			},
		},
	}

	return kb
}

// buttonText returns label with a check or cross prefix
// to indicate selected or not selected.
func buttonText(text string, opt bool) string {
	if opt {
		return "✅ " + text
	}

	return "❌ " + text
}

// StartHandler handles the /start command.
// It sends the keyboard to the chat.
func StartHandler(ctx context.Context, app *bot.Bot, update *models.Update) {
	if update.Message == nil {
		slog.Error("There are no messages")
		return
	}
	kb := EduDegree()
	chatID := update.Message.Chat.ID
	app.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        "Привет! Это твой бот. Select .",
		ReplyMarkup: kb,
	})
}
