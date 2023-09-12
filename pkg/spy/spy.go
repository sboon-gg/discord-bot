package spy

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/bwmarrin/discordgo"
	"github.com/sboon-gg/sby-bot/pkg/config"
	"github.com/sboon-gg/sby-bot/pkg/db"
	"github.com/sboon-gg/sby-bot/pkg/discord"
	"github.com/sboon-gg/sby-bot/pkg/spy/prspy"
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
	bot.refreshRolesCache()
	return bot
}

func (b *Bot) Register(client *discord.Bot) {
	client.RegisterCommand(buttonCommand, b.showButton)
	client.RegisterComponent(infoButton, b.handleButton)
	client.RegisterModal(infoModal, b.handleModal)

	client.RegisterCommand(mapRoleCommand, b.mapRoleHandler)

	go b.roleSetter(client)
}

func (b *Bot) roleSetter(client *discord.Bot) {
	for {
		time.Sleep(time.Minute)

		players := prspy.FetchAllPlayers()
		users := b.userRepo.FindAll()
		b.refreshRolesCache()

		for _, u := range users {
			if p, ok := players[u.IGN]; ok {
				if p.IsAI {
					continue
				}
				b.setActiveRoles(client.S, u.DiscordID)
			} else {
				b.removeActiveRoles(client.S, u.DiscordID)
			}
		}
	}
}

func (b *Bot) refreshRolesCache() {
	roles := b.roleRepo.FindAll()
	b.roleToActiveMap = make(map[string]string)
	for _, role := range roles {
		b.roleToActiveMap[role.DiscordID] = role.ActiveRoleID
		b.activeToRoleMap[role.ActiveRoleID] = role.DiscordID
	}
}

func (b *Bot) setActiveRoles(s *discordgo.Session, discordID string) error {
	member, err := s.GuildMember(b.config.GuildID, discordID)
	if err != nil {
		return err
	}

	for presentRoleID, activeRoleID := range b.roleToActiveMap {
		if presentRoleID == "" {
			// Empty role means @everyone
			err = s.GuildMemberRoleAdd(b.config.GuildID, member.User.ID, activeRoleID)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("unable to set role %s on user %s", activeRoleID, member.User.ID))
			}
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
