package main

import (
	"os"
	"os/signal"
	"syscall"
	"strings"
	"time"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func newUpdateStatusData(idle int, activityType discordgo.ActivityType, name, url string) *discordgo.UpdateStatusData {
	usd := &discordgo.UpdateStatusData{
		Status: "online",
	}

	if idle > 0 {
		usd.IdleSince = &idle
	}

	if name != "" {
		usd.Activities = []*discordgo.Activity{{
			Name: name,
			Type: activityType,
			URL:  url,
		}}
	}

	return usd
}

func BotRun() {
	// create bot session
	dg, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Fatal("error creating session: ", err)
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		log.Fatal("error opening connection: ", err)
	}

        log.Println("Bot is now running.  Press CTRL-C to exit.")

	for true {
		dg.UpdateStatusComplex(*newUpdateStatusData(1, discordgo.ActivityTypeWatching, fmt.Sprintf("%.2f", QuoteDero()) + " USDT", ""))
		time.Sleep(time.Minute)
	}

        // Wait here until CTRL-C or other term signal is received.
        sc := make(chan os.Signal, 1)
        signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
        <-sc

        // Cleanly close down the Discord session.
        dg.Close()
}

func printHelp() (help string) {
	help += "```"
	help += config.Prefix + "help                    this message\n"
	help += config.Prefix + "tokens                  display token info\n"
	help += config.Prefix + "pairs                   display pair info\n"
	help += config.Prefix + "quote <tokenA> <tokenB> price of tokenA in tokenB\n"
//	help += config.Prefix + "channel <channelID>     change bot channel\n"
	help += "```"

	return help
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(config.Channel) > 0 && m.GuildID != "" && m.ChannelID != config.Channel {
		return
	}

	if !strings.HasPrefix(m.Content, config.Prefix) {
		return
	}

	m.Content = strings.TrimPrefix(m.Content, config.Prefix)

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	log.Printf("%s: %s\n", m.Author.String(), m.Content)

	var reply string

	split := strings.Split(m.Content, " ")

	if split[0] == "help" {
		reply = printHelp()
	} else if split[0] == "channel" && m.Author.String() == config.Owner {
		config.Channel = split[1]
		reply = "Changed channel to " + split[1]
	} else if split[0] == "tokens" {
		reply = Tokens()
	} else if split[0] == "pairs" {
		reply = Pairs()
	} else if split[0] == "quote" {
		reply = Quote(split[1:])
	} else {
		reply = "unknown command: " + m.Content
	}

	if len(reply) > 0 {
		s.ChannelMessageSend(m.ChannelID, reply)
	}
}
