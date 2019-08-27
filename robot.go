package main

import (
	"github.com/bwmarrin/discordgo"
	"time"
)

func BuildRobot(session *discordgo.Session, event *discordgo.MessageCreate) {
	UserData := client.Doc("users/" + event.Author.ID)
	_, err = UserData.Set(ctx, User{
		Claim:     time.Now().AddDate(0, 0, -1).String(),
		//TODO: Get rid of this, people can reset their roll timer by just building a new robot allowing other users to farm off of their rolls
		Inventory: nil,
		Mods:      nil,
		Name:      "asdf",
		Rolls:     0,
		Type:      0,
		Victories: 0,
		Defeats:   0,
	})
	return
}

func UpdateRobot(session *discordgo.Session, event *discordgo.MessageCreate) {
	UserData := client.Doc("users/" + event.Author.ID)
	var Inventory []string
	Inventory = append(Inventory, "titaniumgears")
	_, err = UserData.Set(ctx, User{
		Inventory: Inventory,
		//TODO: Make sure to Set any string values, as these will be set to nil when updating data if not provided
	})
	return
}