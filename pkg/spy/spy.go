package spy

import (
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"

	"github.com/bwmarrin/discordgo"
	"github.com/sboon-gg/discord-bot/pkg/config"
	"github.com/sboon-gg/discord-bot/pkg/db"
	"github.com/sboon-gg/discord-bot/pkg/discord"
	"github.com/sboon-gg/discord-bot/pkg/spy/prspy"
)

const (
	everyoneTag = "@everyone"
)

type Bot struct {
	config          *config.Config
	userRepo        *db.UserRepository
	roleRepo        *db.RoleRepository
	roleToActiveMap map[string]string
	activeToRoleMap map[string]string
}

func New(conf *config.Config, userRepo *db.UserRepository, roleRepo *db.RoleRepository) *Bot {
	bot := &Bot{
		config:          conf,
		userRepo:        userRepo,
		roleRepo:        roleRepo,
		roleToActiveMap: make(map[string]string),
		activeToRoleMap: make(map[string]string),
	}
	return bot
}

var spyCommand = &discordgo.ApplicationCommand{
	Name:        "spy",
	Description: "Spy - translate PRSPY activity to Discord role(s)",
	Options: []*discordgo.ApplicationCommandOption{
		mapRoleSubCommand,
		buttonCommand,
	},
}

func (b *Bot) Register(client *discord.Bot) {
	client.RegisterComponent(infoButton, b.handleButton)
	client.RegisterModal(playerInfoModal, b.handleModal)

	client.RegisterCommand(spyCommand, b.commandHandler)

	go func() {
		for {
			err := b.roleSetter(client.Session())
			if err != nil {
				log.Print(err)
			}

			time.Sleep(time.Minute)
		}
	}()
}

func (b *Bot) commandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	switch options[0].Name {
	case buttonCmdName:
		b.showButton(s, i)
	case rolesCmdName:
		b.rolesHandler(s, i)
	}
}

func (b *Bot) roleSetter(s *discordgo.Session) error {
	prspyData, err := prspy.FetchData()
	if err != nil {
		return errors.Wrap(err, "Couldn't fetch PRSPY data")
	}

	players := prspy.GetAllPlayers(prspyData)
	users := b.userRepo.FindAll()
	b.refreshRolesCache(s)

	for _, u := range users {
		if _, ok := players[u.IGN]; ok {
			err = b.setActiveRoles(s, u.DiscordID)
			if err != nil {
				return errors.Wrap(err, "Couldn't set active roles")
			}
		} else {
			err = b.removeActiveRoles(s, u.DiscordID)
			if err != nil {
				return errors.Wrap(err, "Couldn't unset active roles")
			}
		}
	}

	return nil
}

func (b *Bot) refreshRolesCache(s *discordgo.Session) error {
	roles := b.roleRepo.FindAll()

	b.roleToActiveMap = make(map[string]string)
	b.activeToRoleMap = make(map[string]string)

	dRoles, err := s.GuildRoles(b.config.GuildID)
	if err != nil {
		return err
	}

	everyoneID := ""

	for _, d := range dRoles {
		if d.Name == everyoneTag {
			everyoneID = d.ID
		}
	}

	for _, role := range roles {
		id := role.DiscordID
		activeID := role.ActiveRoleID
		if role.DiscordID == everyoneID {
			id = everyoneTag
		}
		b.roleToActiveMap[id] = activeID
		b.activeToRoleMap[activeID] = id
	}

	return nil
}

func (b *Bot) setActiveRoles(s *discordgo.Session, discordID string) error {
	member, err := s.GuildMember(b.config.GuildID, discordID)
	if err != nil {
		return err
	}

	for presentRoleID, activeRoleID := range b.roleToActiveMap {
		if presentRoleID == everyoneTag {
			err = s.GuildMemberRoleAdd(b.config.GuildID, member.User.ID, activeRoleID)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("unable to set role %s on user %s", activeRoleID, member.User.ID))
			}
			continue
		}

		for _, roleID := range member.Roles {
			if presentRoleID == roleID {
				err = s.GuildMemberRoleAdd(b.config.GuildID, member.User.ID, activeRoleID)
				if err != nil {
					return errors.Wrap(err, fmt.Sprintf("unable to set role %s on user %s", activeRoleID, member.User.ID))
				}

				break
			}
		}
	}

	return nil
}

func (b *Bot) removeActiveRoles(s *discordgo.Session, discordID string) error {
	member, err := s.GuildMember(b.config.GuildID, discordID)
	if err != nil {
		return err
	}

	for _, roleID := range member.Roles {
		if _, ok := b.activeToRoleMap[roleID]; ok {
			err = s.GuildMemberRoleRemove(b.config.GuildID, member.User.ID, roleID)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("unable to remove role %s from user %s", roleID, member.User.ID))
			}
		}
	}

	return nil
}
