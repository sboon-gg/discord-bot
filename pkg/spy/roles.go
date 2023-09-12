package spy

import "github.com/bwmarrin/discordgo"

var mapRoleCommand = &discordgo.ApplicationCommand{
	Name:        "maprole",
	Description: "Map discord roles to set when player is active on PRSPY",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionRole,
			Name:        "role",
			Description: "Role user has that will map to active role (@everyone for no restriction)",
			Required:    false,
		},
		{
			Type:        discordgo.ApplicationCommandOptionRole,
			Name:        "active-role",
			Description: "Role user will receive when active in PRSPY (leave empty to remove mapping)",
			Required:    false,
		},
	},
}

func (b *Bot) mapRoleHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {

}
