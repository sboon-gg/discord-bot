package discord

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/sboon-gg/sby-bot/pkg/config"
)

type Handler func(*discordgo.Session, *discordgo.InteractionCreate)

type Command struct {
	Command *discordgo.ApplicationCommand
	Handler Handler
}

type Bot struct {
	Config             *config.Config
	commands           map[string]Command
	components         map[string]Handler
	modals             map[string]Handler
	S                  *discordgo.Session
	registeredCommands []*discordgo.ApplicationCommand
}

func New(conf *config.Config) *Bot {
	s, err := discordgo.New(fmt.Sprintf("Bot %s", conf.Token))
	if err != nil {
		log.Fatal(err)
	}

	c := &Bot{
		Config: conf,
		S:      s,

		commands:           make(map[string]Command),
		components:         make(map[string]Handler),
		modals:             make(map[string]Handler),
		registeredCommands: make([]*discordgo.ApplicationCommand, 0),
	}

	c.prepareSession()

	return c
}

func (c *Bot) prepareSession() {
	c.S.AddHandler(c.interactionsRouter)

	c.S.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
}

func (c *Bot) Run() {
	err := c.S.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	c.createCommands()

	defer c.S.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	c.deleteCommands()
}

func (c *Bot) interactionsRouter(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		if h, ok := c.commands[i.ApplicationCommandData().Name]; ok {
			h.Handler(c.S, i)
		}
	case discordgo.InteractionMessageComponent:
		if h, ok := c.components[i.MessageComponentData().CustomID]; ok {
			h(c.S, i)
		}
	case discordgo.InteractionModalSubmit:
		if h, ok := c.modals[i.ModalSubmitData().CustomID]; ok {
			h(c.S, i)
		}
	}
}

func (c *Bot) RegisterCommand(command *discordgo.ApplicationCommand, handler Handler) {
	c.commands[command.Name] = Command{
		Command: command,
		Handler: handler,
	}
}

func (c *Bot) RegisterComponent(id string, handler Handler) {
	c.components[id] = handler
}

func (c *Bot) RegisterModal(id string, handler Handler) {
	c.modals[id] = handler
}

func (c *Bot) createCommands() {
	for _, command := range c.commands {
		cmd, err := c.S.ApplicationCommandCreate(c.Config.AppID, c.Config.GuildID, command.Command)
		if err != nil {
			log.Fatalf("Cannot create slash command: %v", err)
		}

		c.registeredCommands = append(c.registeredCommands, cmd)
	}
}

func (c *Bot) deleteCommands() {
	log.Println("Removing commands...")
	for _, cmd := range c.registeredCommands {
		err := c.S.ApplicationCommandDelete(c.S.State.User.ID, c.Config.GuildID, cmd.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", cmd.Name, err)
		}
	}
}
