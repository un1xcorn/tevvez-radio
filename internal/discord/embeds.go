package discord

import (
	"github.com/bwmarrin/discordgo"
)

func (d *Discord) PrefixEmbed() *discordgo.MessageEmbed {
	embed := discordgo.MessageEmbed{}
	embed.Title = "My prefix is: `" + d.Prefix + "`"
	embed.Color = d.Color

	return &embed
}

func (d *Discord) ErrorEmbed() *discordgo.MessageEmbed {
	embed := discordgo.MessageEmbed{}
	embed.Title = "An error happened."
	embed.Color = 16711680

	return &embed
}
