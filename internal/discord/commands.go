package discord

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Name        string
	Description string
	Usage       string
	Category    string
	Aliases     []string
	Args        int
	Perms       []int64
	Run         func(s *discordgo.Session, m *discordgo.MessageCreate, args []string)
}

func (d *Discord) LoadCommands() {
	rawCommands := []Command{}

	rawCommands = append(rawCommands, Command{
		Name:        "help",
		Description: "List the available commands.",
		Usage:       "help",
		Category:    "Generalℹ️",
		Aliases:     []string{"h", "commands", "command", "how"},
		Args:        0,
		Perms:       nil,
		Run: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			embed := discordgo.MessageEmbed{}
			embed.Title = "Commands list"
			embed.Color = d.Color
			embed.Description += "`usage` `(aliases)`: *description*\n\n"

			for categoryName, categoryCommands := range d.Categories {
				embed.Description += "__" + categoryName + "__\n"
				for _, commandName := range categoryCommands {
					command := d.Commands[commandName]
					embed.Description += "`" + d.Prefix + command.Usage + "`"
					embed.Description += " `(" + strings.Join(command.Aliases, " | ") + ")`: "
					embed.Description += "*" + command.Description + "*\n"
				}
				embed.Description += "\n"
			}

			s.ChannelMessageSendEmbed(m.ChannelID, &embed)
		}})

	rawCommands = append(rawCommands, Command{
		Name:        "play",
		Description: "Start playing the radio.",
		Usage:       "play",
		Category:    "Music🎵",
		Aliases:     []string{},
		Args:        0,
		Perms:       nil,
		Run: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			memberVoiceChatID := GetMemberVoiceChatID(s, m.GuildID, m.Author.ID)

			if memberVoiceChatID != "" && !d.MusicPlayer.IsPlayerPlaying(m.GuildID) {
				embed := discordgo.MessageEmbed{}
				embed.Title = "Started playing the radio."
				embed.Color = d.Color

				s.ChannelMessageSendEmbed(m.ChannelID, &embed)
				d.MusicPlayer.PlayRandomTrack(m.GuildID, memberVoiceChatID)
			} else {
				s.ChannelMessageSendEmbed(m.ChannelID, d.ErrorEmbed())
			}
		}})

	rawCommands = append(rawCommands, Command{
		Name:        "stop",
		Description: "Stop playing the radio.",
		Usage:       "stop",
		Category:    "Music🎵",
		Aliases:     []string{},
		Args:        0,
		Perms:       nil,
		Run: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			memberVoiceChatID := GetMemberVoiceChatID(s, m.GuildID, m.Author.ID)
			botVoiceChatID := GetMemberVoiceChatID(s, m.GuildID, s.State.User.ID)

			if memberVoiceChatID == botVoiceChatID && d.MusicPlayer.IsPlayerPlaying(m.GuildID) {
				embed := discordgo.MessageEmbed{}
				embed.Title = "Stopped playing the radio."
				embed.Color = d.Color

				s.ChannelMessageSendEmbed(m.ChannelID, &embed)
				d.MusicPlayer.RemovePlayer(m.GuildID)
			} else {
				s.ChannelMessageSendEmbed(m.ChannelID, d.ErrorEmbed())
			}
		}})

	rawCommands = append(rawCommands, Command{
		Name:        "current",
		Description: "Get the current song.",
		Usage:       "current",
		Category:    "Music🎵",
		Aliases:     []string{},
		Args:        0,
		Perms:       nil,
		Run: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			if d.MusicPlayer.IsPlayerPlaying(m.GuildID) {
				infos := d.MusicPlayer.GetCurrentTrack(m.GuildID)
				embed := discordgo.MessageEmbed{}
				embed.Title = "Current track"
				embed.Color = d.Color
				embed.URL = *infos.URI
				embed.Description = "Name: **" + infos.Title + "**\nLength: **" + infos.Length.String() + "**"
				embed.Image = &discordgo.MessageEmbedImage{
					URL: "https://i.ytimg.com/vi/" + infos.Identifier + "/hqdefault.jpg",
				}

				s.ChannelMessageSendEmbed(m.ChannelID, &embed)
			} else {
				s.ChannelMessageSendEmbed(m.ChannelID, d.ErrorEmbed())
			}
		}})

	commands := make(map[string]Command)
	aliases := make(map[string]string)
	categories := make(map[string][]string)

	for _, command := range rawCommands {
		commands[command.Name] = command
		for _, alias := range command.Aliases {
			aliases[alias] = command.Name
		}
	}
	for _, command := range commands {
		categories[command.Category] = append(categories[command.Category], command.Name)
	}

	d.Commands = commands
	d.Aliases = aliases
	d.Categories = categories
}
