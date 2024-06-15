package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func main() {
	slog.Info("initiate helath check routine...")
	go func() {
		_ = http.ListenAndServe(":8080", http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("alive"))
			},
		))
	}()

	slog.Info("starting w15p...")
	slog.Info("disgo version", slog.String("version", disgo.Version))

	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		slog.Error("DISCORD_BOT_TOKEN has not been set.")
		return
	}

	client, err := disgo.New(token,
		// set gateway options
		bot.WithGatewayConfigOpts(
			// set enabled intents
			gateway.WithIntents(
				gateway.IntentGuildMessages,
				gateway.IntentMessageContent,
			),
		),
		// add event listeners
		bot.WithEventListenerFunc(onMessageCreate),
	)

	if err != nil {
		slog.Error("error while building disgo", slog.Any("err", err))
		return
	}

	defer client.Close(context.TODO())

	// connect to the gateway
	if err = client.OpenGateway(context.TODO()); err != nil {
		slog.Error("errors while connecting to gateway", slog.Any("err", err))
		return
	}

	slog.Info("w15p is awake and ready to interact")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s
}

func onMessageCreate(event *events.MessageCreate) {
	if event.Message.Author.Bot {
		return
	}
	var message string
	if event.Message.Content == "beep" {
		message = "boop"
	} else if event.Message.Content == "boop" {
		message = ":elfgun:"
	}
	if message != "" {
		_, _ = event.Client().Rest().CreateMessage(event.ChannelID, discord.NewMessageCreateBuilder().SetContent(message).Build())
	}
}
