package main

import (
	"flag"

	"github.com/sboon-gg/sboon-bot/pkg/config"
	"github.com/sboon-gg/sboon-bot/pkg/db"
	"github.com/sboon-gg/sboon-bot/pkg/discord"
	"github.com/sboon-gg/sboon-bot/pkg/spy"
)

func main() {
	configFileName := flag.String("config", "config.yaml", "Config file path")
	flag.Parse()

	conf, err := config.New(*configFileName)
	if err != nil {
		panic(err)
	}

	conn := db.New(&conf.Db)
	userRepo := db.NewUserRepository(conn)
	roleRepo := db.NewRoleRepository(conn)

	bot := discord.New(conf)

	spyBot := spy.New(conf, userRepo, roleRepo)
	spyBot.Register(bot)

	bot.Run()
}
