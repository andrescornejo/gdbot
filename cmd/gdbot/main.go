package main

import (
	"log/slog"
	"regexp"

	//"strconv"

	"gdbot/internal/gdbot"
)

var (
	urlPattern    = regexp.MustCompile("^https?://[-a-zA-Z0-9+&@#/%?=~_|!:,.;]*[-a-zA-Z0-9+&@#/%=~_|]?")
	searchPattern = regexp.MustCompile(`^(.{2})search:(.+)`)
)

func main() {
	slog.Info("Starting GDbot...")
	slog.SetLogLoggerLevel(slog.LevelDebug)

	cfg, err := gdbot.LoadConfig()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	b := gdbot.New(*cfg)

	b.StartAndBlock()
	// registerCommands(client)
	//
	// b.Handlers = map[string]func(event *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) error{
	// 	"play":        b.play,
	// 	"pause":       b.pause,
	// 	"now-playing": b.nowPlaying,
	// 	"stop":        b.stop,
	// 	"players":     b.players,
	// 	"queue":       b.queue,
	// 	"clear-queue": b.clearQueue,
	// 	"queue-type":  b.queueType,
	// 	"shuffle":     b.shuffle,
	// 	"seek":        b.seek,
	// 	"volume":      b.volume,
	// 	"skip":        b.skip,
	// 	"bass-boost":  b.bassBoost,
	// }
}
