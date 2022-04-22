package discord

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (d *Discord) ReadyEvent(session *discordgo.Session, r *discordgo.Ready) {
	session.UpdateGameStatus(0, fmt.Sprint(len(d.MusicPlayer.Tracks))+" tracks.")
}

func (d *Discord) MessageCreateEvent(session *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == session.State.User.ID {
		return
	}

	if IDTagged(session.State.User.ID, m.Mentions) {
		session.ChannelMessageSendEmbed(m.ChannelID, d.PrefixEmbed())
		return
	}

	if !strings.HasPrefix(m.Content, d.Prefix) {
		return
	}

	permissions, _ := session.UserChannelPermissions(m.Author.ID, m.ChannelID)
	if command, ok := d.Commands[GetCommandName(m.Content, d.Prefix)]; ok {
		if UserHasCommandPerms(command, permissions) {
			if len(GetCommandArgs(m.Content)) >= command.Args {
				command.Run(session, m, GetCommandArgs(m.Content))
			} else {
				session.ChannelMessageSendEmbed(m.ChannelID, d.ErrorEmbed())
			}
		} else {
			session.ChannelMessageSendEmbed(m.ChannelID, d.ErrorEmbed())
		}
	}
	if command, ok := d.Commands[d.Aliases[GetCommandName(m.Content, d.Prefix)]]; ok {
		if UserHasCommandPerms(command, permissions) {
			if len(GetCommandArgs(m.Content)) >= command.Args {
				command.Run(session, m, GetCommandArgs(m.Content))
			} else {
				session.ChannelMessageSendEmbed(m.ChannelID, d.ErrorEmbed())
			}
		} else {
			session.ChannelMessageSendEmbed(m.ChannelID, d.ErrorEmbed())
		}
	}
}
