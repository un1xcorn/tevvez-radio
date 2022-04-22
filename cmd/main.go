package main

import (
	"github.com/un1xcorn/tevvez-radio/internal/discord"
	"github.com/un1xcorn/tevvez-radio/pkg/config"
	"github.com/un1xcorn/tevvez-radio/pkg/database"
)

func main() {
	config, _ := config.Load("config.json")
	database, _ := database.GetDatabase("database.sqlite")
	defer database.Close()

	bot := discord.Discord{
		Token:    config.Token,
		Color:    config.Color,
		Prefix:   config.Prefix,
		Database: database,
	}
	bot.Run()
}
