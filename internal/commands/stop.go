package commands

import (
	"context"
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"

	"github.com/disgoorg/disgolink/v3/lavalink"

	"gdbot/internal/gdbot"
	"gdbot/utils"
)

var stopCommand = discord.SlashCommandCreate{
	Name:        "stop",
	Description: "Stops the current song and stops the player",
}

// func  stop(b *gdbot.GDBot) (event *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) error {
// 	player := b.Lavalink.ExistingPlayer(*event.GuildID())
// 	if player == nil {
// 		return event.CreateMessage(discord.MessageCreate{
// 			Content: "No player found",
// 		})
// 	}
//
// 	if err := player.Update(context.TODO(), lavalink.WithNullTrack()); err != nil {
// 		return event.CreateMessage(discord.MessageCreate{
// 			Content: fmt.Sprintf("Error while stopping: `%s`", err),
// 		})
// 	}
//
// 	return event.CreateMessage(discord.MessageCreate{
// 		Content: "Player stopped",
// 	})
// }

func HandleStop(b *gdbot.GDBot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		player := b.Lavalink.ExistingPlayer(*e.GuildID())
		if player == nil {
			return utils.CreateBasicMessage(e.Respond, "Error", "No player found", utils.ColorError)
		}

		if err := player.Update(context.TODO(), lavalink.WithNullTrack()); err != nil {
			err_msg := fmt.Sprintf("Error while stopping: `%s`", err)

			return utils.CreateBasicMessage(e.Respond, "Error", err_msg, utils.ColorError)
		}

		return utils.CreateBasicMessage(e.Respond, "Stop", "Player stopped", utils.ColorSuccess)
	}
}
