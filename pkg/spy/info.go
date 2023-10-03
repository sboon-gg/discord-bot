package spy

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

const (
	buttonCmdName   = "button"
	infoButton      = "infoButton"
	playerInfoModal = "infoModal"
)

var buttonCommand = &discordgo.ApplicationCommandOption{
	Name:        buttonCmdName,
	Description: "Display button for users to put their info in",
	Type:        discordgo.ApplicationCommandOptionSubCommand,
}

var button = discordgo.Button{
	Label:    "Player info form",
	Style:    discordgo.SuccessButton,
	CustomID: infoButton,
}

var buttonResponse = discordgo.InteractionResponse{
	Type: discordgo.InteractionResponseChannelMessageWithSource,
	Data: &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title: "PR Activity",
				Description: `
By activating this bot, whenever you are playing PR on a public server, you will be assigned a role on SBOON Discord to inform that you are in-game.
To activate the bot, just insert your in-game name (without tag).

**Click** the **button** below to **start**.`,
			},
		},
		Components: []discordgo.MessageComponent{
			&discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					button,
				},
			},
		},
	},
}

func (b *Spy) showButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &buttonResponse)
	if err != nil {
		log.Print(err)
	}
}

func (b *Spy) handleButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	user := b.userRepo.FindByDiscordID(i.Member.User.ID)
	ign := ""
	// hash := ""
	if user != nil {
		ign = user.IGN
		// hash = user.Hash
	}

	modal := discordgo.InteractionResponseData{
		CustomID: playerInfoModal,
		Title:    "Player info",
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "ign",
						Label:       "What is your in-game name (without tag)?",
						Style:       discordgo.TextInputShort,
						Placeholder: "",
						Required:    true,
						MaxLength:   20,
						MinLength:   4,
						Value:       ign,
					},
				},
			},
			// discordgo.ActionsRow{
			// 	Components: []discordgo.MessageComponent{
			// 		discordgo.TextInput{
			// 			CustomID:    "hash",
			// 			Label:       "What is your hash?",
			// 			Style:       discordgo.TextInputShort,
			// 			Placeholder: "",
			// 			Required:    false,
			// 			MaxLength:   32,
			// 			MinLength:   32,
			// 			Value:       hash,
			// 		},
			// 	},
			// },
		},
	}

	modalResponse := discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &modal,
	}

	err := s.InteractionRespond(i.Interaction, &modalResponse)
	if err != nil {
		log.Print(err)
	}
}

func (b *Spy) handleModal(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ModalSubmitData()
	ign := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	// hash := data.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	b.userRepo.SetInfo(i.Member.User.ID, ign, "")

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: `
Thank you.

Next time you get on the battlefield, you will be assigned the "Playing" role.
Use this to coordinate with your mates, and have fun!
`,
			Embeds: []*discordgo.MessageEmbed{
				{
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:  "In-game name",
							Value: ign,
						},
						// {
						// 	Name:  "Hash",
						// 	Value: hash,
						// },
					},
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Print(err)
	}
}
