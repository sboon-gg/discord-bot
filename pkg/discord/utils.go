package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func ErrorResponse(s *discordgo.Session, i *discordgo.Interaction, err error) {
	respErr := s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					// Color: discordgo.
					Fields: []*discordgo.MessageEmbedField{
						{
							Value: err.Error(),
						},
					},
				},
			},
		},
	})

	if respErr != nil {
		log.Print(respErr)
	}
}
