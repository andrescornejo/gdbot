// Package gdbot TODO
//
// It includes functions like Add, Subtract, and Multiply for everyday arithmetic,
// with support for both integers and floats. See the individual function docs
// for usage examples.
package gdbot

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"

	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/snowflake/v2"
)

type GDBot struct {
	Client   bot.Client
	Lavalink disgolink.Client
	Handlers map[string]func(event *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) error
	Queues   *QueueManager
}

func New(cfg Config) *GDBot {
	b := &GDBot{
		Queues: &QueueManager{
			queues: make(map[snowflake.ID]*Queue),
		},
	}
	b.initGDBotServices(cfg)
	return b
}

func (b *GDBot) StartAndBlock() {
	err := b.Client.OpenGateway(context.TODO())
	if err != nil {
		slog.Error("failed to open gateway", slog.Any("err", err))
		os.Exit(1)
	}

	slog.Info("Connected to Discord")
	defer b.Client.Close(context.TODO())

	slog.Info("GDBot is running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

func (b *GDBot) initGDBotServices(cfg Config) {
	b.Client = b.buildDisgoClient(cfg.Token)
	b.Lavalink = b.buildLavalinkClient(b.Client.ApplicationID(), context.TODO(), cfg.NodeName, cfg.NodeAddress, cfg.NodePassword, cfg.NodeSecure)
}

func (b *GDBot) buildDisgoClient(token string) bot.Client {
	client, err := disgo.New(token,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(gateway.IntentGuilds, gateway.IntentGuildVoiceStates),
		),
		bot.WithCacheConfigOpts(
			cache.WithCaches(cache.FlagVoiceStates),
		),
		bot.WithEventListenerFunc(b.onApplicationCommand),
		bot.WithEventListenerFunc(b.onVoiceStateUpdate),
		bot.WithEventListenerFunc(b.onVoiceServerUpdate),
	)
	if err != nil {
		slog.Error("error while building disgo client", slog.Any("err", err))
		os.Exit(1)
	}
	return client
}

func (b *GDBot) buildLavalinkClient(appID snowflake.ID, ctx context.Context, nodeName string, nodeAddress string, nodePassword string, nodeSecure bool) disgolink.Client {
	lavalink := disgolink.New(appID,
		disgolink.WithListenerFunc(b.onPlayerPause),
		disgolink.WithListenerFunc(b.onPlayerResume),
		disgolink.WithListenerFunc(b.onTrackStart),
		disgolink.WithListenerFunc(b.onTrackEnd),
		disgolink.WithListenerFunc(b.onTrackException),
		disgolink.WithListenerFunc(b.onTrackStuck),
		disgolink.WithListenerFunc(b.onWebSocketClosed),
		disgolink.WithListenerFunc(b.onUnknownEvent),
	)

	node, err := lavalink.AddNode(ctx, disgolink.NodeConfig{
		Name:     nodeName,
		Address:  nodeAddress,
		Password: nodePassword,
		Secure:   nodeSecure,
	})
	if err != nil {
		slog.Error("failure adding lavalink node", slog.Any("err", err))
		os.Exit(1)
	}
	version, err := node.Version(ctx)
	if err != nil {
		slog.Error("failure getting lavalink node version", slog.Any("err", err))
		os.Exit(1)
	}

	slog.Info("Utilizing Lavalink version:", slog.String("version", version))
	slog.Info("Lavalink node session ID:", slog.String("sessionID", node.SessionID()))

	return lavalink
}

func (b *GDBot) syncCommands(commands []discord.ApplicationCommandCreate, guildIDs ...snowflake.ID) {
	if len(guildIDs) == 0 {
		if _, err := b.Client.Rest().SetGlobalCommands(b.Client.ApplicationID(), commands); err != nil {
			slog.Error("Failed to sync commands: %s", slog.Any("error", err))
		}
		return
	}
	for _, guildID := range guildIDs {
		if _, err := b.Client.Rest().SetGuildCommands(b.Client.ApplicationID(), guildID, commands); err != nil {
			slog.Error("Failed to sync commands: %s", slog.Any("error", err))
		}
	}
}

func (b *GDBot) onApplicationCommand(event *events.ApplicationCommandInteractionCreate) {
	data := event.SlashCommandInteractionData()

	handler, ok := b.Handlers[data.CommandName()]
	if !ok {
		slog.Info("unknown command", slog.String("command", data.CommandName()))
		return
	}
	if err := handler(event, data); err != nil {
		slog.Error("error handling command", slog.Any("err", err))
	}
}

func (b *GDBot) onVoiceStateUpdate(event *events.GuildVoiceStateUpdate) {
	if event.VoiceState.UserID != b.Client.ApplicationID() {
		return
	}
	b.Lavalink.OnVoiceStateUpdate(context.TODO(), event.VoiceState.GuildID, event.VoiceState.ChannelID, event.VoiceState.SessionID)
	if event.VoiceState.ChannelID == nil {
		b.Queues.Delete(event.VoiceState.GuildID)
	}
}

func (b *GDBot) onVoiceServerUpdate(event *events.VoiceServerUpdate) {
	b.Lavalink.OnVoiceServerUpdate(context.TODO(), event.GuildID, event.Token, *event.Endpoint)
}
