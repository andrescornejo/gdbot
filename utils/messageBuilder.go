package utils

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func CreateBasicMessage(respondFunc events.InteractionResponderFunc, title string, content string, color int) error {
	return respondFunc(discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageCreateBuilder().
			SetEmbeds(discord.NewEmbedBuilder().
				SetTitle(title).
				SetColor(color).
				AddField("Content", content, false).
				Build(),
			))
}
