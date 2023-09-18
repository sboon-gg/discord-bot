package main

import (
	"github.com/sboon-gg/sboon-bot/pkg/config"
	"github.com/sboon-gg/sboon-bot/pkg/db"
	"github.com/sboon-gg/sboon-bot/pkg/discord"
	"github.com/sboon-gg/sboon-bot/pkg/spy"
)

func main() {
	conf, err := config.New()
	if err != nil {
		panic(err)
	}

	conn := db.New()
	userRepo := db.NewUserRepository(conn)
	roleRepo := db.NewRoleRepository(conn)

	bot := discord.New(conf)

	spyBot := spy.New(conf, userRepo, roleRepo)
	spyBot.Register(bot)

	bot.Run()
}
