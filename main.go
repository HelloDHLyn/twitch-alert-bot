package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"./app"
	"./twitch"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron"
)

var (
	channelId string
)

func main() {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Add cron job to check game changes.
	c := cron.New()
	c.AddFunc("*/30 * * * * *", func() { checkGameChanges(dg) })
	c.Start()

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	commands := strings.Split(m.Content, " ")
	switch commands[0] {
	// !스트리머 목록
	// !스트리머 추가 <LoginID>
	// !스트리머 삭제 <LoginID>
	case "!스트리머":
		if len(commands) < 2 {
			s.ChannelMessageSend(m.ChannelID, "!스트리머 목록\n!스트리머 추가 <LoginID>\n!스트리머 삭제 <LoginID>")
			return
		}

		switch commands[1] {
		case "목록":
			users := twitch.GetUsers(app.GetFollowedUsers()).Data
			msg := ""
			for _, data := range users {
				msg = msg + data.DisplayName + "(" + data.Login + ")\n"
			}
			s.ChannelMessageSend(m.ChannelID, msg)

		case "추가":
			users := twitch.GetUsers([]string{commands[2]}).Data
			if len(users) == 0 {
				s.ChannelMessageSend(m.ChannelID, "해당하는 ID의 스트리머가 없습니다.")
				return
			}
			app.AddFollowedUser(commands[2])
			s.ChannelMessageSend(m.ChannelID, users[0].DisplayName+"("+commands[2]+")"+" 스트리머를 팔로우했습니다 :)")

		case "삭제":
			users := twitch.GetUsers([]string{commands[2]}).Data
			if len(users) == 0 {
				s.ChannelMessageSend(m.ChannelID, "해당하는 ID의 스트리머가 없습니다.")
				return
			}
			app.DeleteFollowedUser(commands[2])
			s.ChannelMessageSend(m.ChannelID, users[0].DisplayName+"("+commands[2]+")"+" 스트리머를 언팔로우했습니다 :(")
		}

	// !봇하
	case "!봇하":
		channelId = m.ChannelID
		s.ChannelMessageSend(m.ChannelID, "<:ramelHihi:463640108562907138><:ramelHihi:463640108562907138><:ramelHihi:463640108562907138>")

	case "!봇바":
		channelId = ""

	// !명령어
	case "!명령어":
		s.ChannelMessageSend(m.ChannelID, "!스트리머, !봇하, !명령어")
	}
}

func checkGameChanges(s *discordgo.Session) {
	if channelId == "" {
		return
	}

	for _, user := range twitch.GetUsers(app.GetFollowedUsers()).Data {
		for _, stream := range twitch.GetStreams([]string{user.Login}).Data {
			lastPlayedGameId := app.GetLastPlayedGameId(user.Login)
			if lastPlayedGameId == stream.GameID {
				continue
			}

			game := twitch.GetGames([]string{stream.GameID}).Data[0]
			s.ChannelMessageSendComplex(channelId, &discordgo.MessageSend{
				Content: user.DisplayName + " - 플레이중인 게임이 변경되었습니다.",
				Embed: &discordgo.MessageEmbed{
					Title:       "플레이 중인 게임",
					Description: fmt.Sprintf("%s", game.Name),
					URL:         fmt.Sprintf("https://twitch.tv/%s", user.Login),
					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL: user.ProfileImageURL,
					},
				},
			})

			app.UpdateLastPlayedGameId(user.Login, game.ID)
		}
	}
}
