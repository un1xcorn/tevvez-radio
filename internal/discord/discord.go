package discord

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/un1xcorn/tevvez-radio/internal/musicplayer"
	"github.com/un1xcorn/tevvez-radio/pkg/database"
)

type Discord struct {
	Token       string
	Color       int
	Prefix      string
	Commands    map[string]Command
	Aliases     map[string]string
	Categories  map[string][]string
	Database    database.Database
	MusicPlayer *musicplayer.MusicPlayer
}

func (d *Discord) Run() {
	session, err := discordgo.New("Bot " + d.Token)
	if err != nil {
		log.Fatal("An error happened.")
	}

	session.AddHandler(d.ReadyEvent)
	session.AddHandler(d.MessageCreateEvent)

	d.LoadCommands()
	log.Info(len(d.Commands), " commands have been loaded.")

	tracks, _ := d.Database.GetTracks()
	d.MusicPlayer = musicplayer.New(session, tracks)

	err = session.Open()
	if err != nil {
		log.Fatal("An error happened.")
	}

	log.Info("Bot is now running.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)
	<-sc
	session.Close()
	log.Info("Bot has been stopped.")
}

func GetCommandName(command string, prefix string) string {
	return strings.Split(strings.ReplaceAll(command, prefix, ""), " ")[0]
}

func GetCommandArgs(command string) []string {
	return strings.Split(command, " ")[1:]
}

func UserHasCommandPerms(command Command, permissions int64) bool {
	for _, permission := range command.Perms {
		if permissions&permission != permission {
			return false
		}
	}
	return true
}

func IDTagged(id string, mentions []*discordgo.User) bool {
	for _, mention := range mentions {
		if mention.ID == id {
			return true
		}
	}
	return false
}

func GetMemberVoiceChatID(session *discordgo.Session, guildID string, id string) string {
	guild, _ := session.State.Guild(guildID)
	for _, voiceState := range guild.VoiceStates {
		if voiceState.UserID == id {
			return voiceState.ChannelID
		}
	}
	return ""
}
