package musicplayer

import (
	"context"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/disgolink/dgolink"
	"github.com/disgoorg/disgolink/lavalink"
	"github.com/disgoorg/snowflake"
)

type MusicPlayer struct {
	lavalink.PlayerEventAdapter
	Session  *discordgo.Session
	Instance *dgolink.Link
	Players  map[string]lavalink.Player
	Tracks   []string
}

func New(session *discordgo.Session, tracks []string) *MusicPlayer {
	musicPlayer := &MusicPlayer{
		Session:  session,
		Instance: dgolink.New(session),
		Players:  make(map[string]lavalink.Player),
		Tracks:   tracks,
	}
	musicPlayer.Instance.AddNode(context.TODO(), lavalink.NodeConfig{
		Name:        "base",
		Host:        "localhost",
		Port:        "1337",
		Password:    "h4ck1ngb4s3d",
		Secure:      false,
		ResumingKey: "",
	})

	return musicPlayer
}

func (m *MusicPlayer) PlayRandomTrack(guildID string, voiceID string) {
	rand.Seed(time.Now().Unix())
	track := m.Tracks[rand.Intn(len(m.Tracks))]

	m.Instance.BestRestClient().LoadItemHandler(context.TODO(), track, lavalink.NewResultHandler(
		func(track lavalink.AudioTrack) {
			m.Session.ChannelVoiceJoinManual(guildID, voiceID, false, false)
			m.GetPlayer(guildID).Play(track)
		},
		func(playlist lavalink.AudioPlaylist) {},
		func(tracks []lavalink.AudioTrack) {},
		func() {},
		func(ex lavalink.FriendlyException) {},
	))
}

func (m *MusicPlayer) IsPlayerPlaying(guildID string) bool {
	return m.Players[guildID] != nil
}

func (m *MusicPlayer) GetCurrentTrack(guildID string) lavalink.AudioTrackInfo {
	return m.Players[guildID].PlayingTrack().Info()
}

func (m *MusicPlayer) GetPlayer(guildID string) lavalink.Player {
	if !m.IsPlayerPlaying(guildID) {
		m.Players[guildID] = m.Instance.Player(snowflake.Snowflake(guildID))
		m.Players[guildID].AddListener(m)
	}

	return m.Players[guildID]
}

func (m *MusicPlayer) RemovePlayer(guildID string) {
	m.Session.ChannelVoiceJoinManual(guildID, "", false, false)
	m.Players[guildID].Destroy()
	delete(m.Players, guildID)
}

func (m *MusicPlayer) OnTrackEnd(player lavalink.Player, track lavalink.AudioTrack, endReason lavalink.AudioTrackEndReason) {
	m.PlayRandomTrack(player.GuildID().String(), player.ChannelID().String())
}
