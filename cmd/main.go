// Entry point for the Telegram bot.
// Sets up logger and environment, creates the bot,
// registers handlers, and starts polling updates.
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/joho/godotenv"
	internal "github.com/one-zero-eight/broadcast/internal/handlers"
)

// main initializes logger and env, starts the bot,
// registers handlers, and blocks until Ctrl+C.
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	// Create text logger to stdout.
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Load environment variables from .env (if present).
	err := godotenv.Load(".env")
	if err != nil {
		slog.Error("Error loading .env file")
	}

	// Read Telegram token from environment.
	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		// If token is missing, log and exit.
		slog.Info("TELEGRAM_TOKEN environment variable not set")
		return
	}

	// Create Telegram bot client.
	app, err := bot.New(token)
	if err != nil {
		// Log error and stop if bot cannot start.
		slog.Error("Error while starting the application", err)
		return
	} else {
		slog.Info("Bot is now running.  Press CTRL-C to exit.")
	}

	// Register /start command handler.
	app.RegisterHandler(bot.HandlerTypeMessageText, "start", bot.MatchTypeCommand, internal.StartHandler)
	// Register callback handler for any data starting with "btn_".
	app.RegisterHandler(bot.HandlerTypeCallbackQueryData, "btn_", bot.MatchTypePrefix, internal.CallbackHandler)

	// Start polling updates and block until context is canceled.
	app.Start(ctx)

}
