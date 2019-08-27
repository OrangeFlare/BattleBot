package main

import (
	"flag"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"cloud.google.com/go/firestore"
	"github.com/bwmarrin/discordgo"
)

type User struct {
	Claim		string		`firestore:"claim"`
	Inventory	[]string	`firestore:"inventory"`
	Mods		[]string	`firestore:"modslot"`
	Name		string		`firestore:"name"`
	Rolls		float64		`firestore:"rolls"`
	Type		float64		`firestore:"type"`
	Victories	float64		`firestore:"victories"`
	Defeats		float64		`firestore:"defeats"`
}

func init() {
	flag.StringVar(&DiscordToken, "t", "", "Discord API Token")
	flag.StringVar(&CredentialsFile, "c", "", "GCP Credentials File")
	flag.StringVar(&OwnerID, "o", "176108182056206336", "Discord User ID")
	flag.Parse()
}

var	DiscordToken	string
var CredentialsFile	string
var	OwnerID			string
var	client			*firestore.Client
var	err				error
var	ctx				context.Context

func main() {
	if DiscordToken == "" {
		fmt.Println("==BattleBot==\n[Error] Your start command is invalid!\nPlease provide -t <Discord API Token>")
		os.Exit(0)
	}

	ctx = context.Background()
	client, err = firestore.NewClient(ctx, "battlebot-250917", option.WithCredentialsFile(CredentialsFile))
	if err != nil {
		fmt.Println("==BattleBot==\n[Fatal Error] " + err.Error())
		os.Exit(0)
	}

	bb, err := discordgo.New("Bot " + DiscordToken)
	if err != nil {
		fmt.Println("==BattleBot==\n[Fatal Error] " + err.Error())
		os.Exit(0)
	}

	bb.AddHandler(BotReady)
	bb.AddHandler(MessageCreateHandler)

	err = bb.Open()
	if err != nil {
		fmt.Println("==BattleBot==\n[Error] " + err.Error())
	}

	fmt.Println("==BattleBot==\n[Info] BattleBot Online and Ready to Fight!")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	err = bb.Close()
	if err != nil {
		fmt.Println("==BattleBot==\n[Error] " + err.Error())
	}
	err = client.Close()
	if err != nil {
		fmt.Println("==BattleBot==\n[Error] " + err.Error())
	}
	return
}

func BotReady(session *discordgo.Session, event *discordgo.Ready) {
	err := session.UpdateStatus(0, "Charging batteries ...")
	if err != nil {
		fmt.Println("==BattleBot==\n[Error] " + err.Error())
	}
	return
}

func MessageCreateHandler(session *discordgo.Session, event *discordgo.MessageCreate) {
	if match, _ := regexp.MatchString("(?i)^bb.reboot$", event.Content); match == true && event.Author.ID == OwnerID {
		fmt.Println("[Info] Restarting Bot ...")
		err := session.ChannelMessageDelete(event.ChannelID, event.ID)
		if err != nil {
			fmt.Println("==BattleBot==\n[Error] " + err.Error())
		}
		err = session.Close()
		if err != nil {
			fmt.Println("==BattleBot==\n[Error] " + err.Error())
		}
		err = client.Close()
		if err != nil {
			fmt.Println("==BattleBot==\n[Error] " + err.Error())
		}
		go main()
		return
	}
	if match, _ := regexp.MatchString("(?i)^bb.buildrobot$", event.Content); match == true {
		go BuildRobot(session, event)
		return
	}
	if match, _ := regexp.MatchString("(?i)^bb.updatearray$", event.Content); match == true {
		go UpdateRobot(session, event)
		return
	}
	return
}