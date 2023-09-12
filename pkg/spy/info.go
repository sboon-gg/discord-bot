package spy

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

const (
	infoButton = "infoButton"
	infoModal  = "infoModal"
)

var buttonCommand = &discordgo.ApplicationCommand{
	Name:        "button",
	Description: "Test the buttons if you got courage",
}

var button = discordgo.Button{
	Label:    "Insert info",
	Style:    discordgo.SuccessButton,
	CustomID: infoButton,
}

var buttonResponse = discordgo.InteractionResponse{
	Type: discordgo.InteractionResponseChannelMessageWithSource,
	Data: &discordgo.InteractionResponseData{
		Components: []discordgo.MessageComponent{
			&discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					button,
				},
			},
		},
	},
}

func (b *Bot) showButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &buttonResponse)
	if err != nil {
		log.Fatal(err)
	}
}

func (b *Bot) handleButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// TODO: read existing from DB
	err := s.InteractionRespond(i.Interaction, &modalResponse)
	if err != nil {
		log.Fatal(err)
	}
}

var modalResponse = discordgo.InteractionResponse{
	Type: discordgo.InteractionResponseModal,
	Data: &modal,
}

var modal = discordgo.InteractionResponseData{
	CustomID: infoModal,
	Title:    "Info modal",
	Components: []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.TextInput{
					CustomID:    "ign",
					Label:       "What is your in-game name (without tag)?",
					Style:       discordgo.TextInputShort,
					Placeholder: "cassius23",
					Required:    true,
					MaxLength:   20,
					MinLength:   4,
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.TextInput{
					CustomID:    "hash",
					Label:       "What is your hash?",
					Style:       discordgo.TextInputShort,
					Placeholder: "",
					Required:    false,
					MaxLength:   32,
					MinLength:   32,
				},
			},
		},
	},
}

func (b *Bot) handleModal(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ModalSubmitData()
	ign := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	hash := data.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	// TODO: save into DB

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Thank you",
			Embeds: []*discordgo.MessageEmbed{
				{
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:  "In-game name",
							Value: ign,
						},
						{
							Name:  "Hash",
							Value: hash,
						},
					},
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
