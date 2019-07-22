package main

import (
	"github.com/bwmarrin/discordgo"
)

func CreateNewMenu(MenuID string, msg discordgo.Message) *EmbedMenu {
	MainMenu := NewEmbedMenu(
		discordgo.MessageEmbed{
			Title: "Freethewump.us Manager",
			Description: "Using this bot, you can manage your domains and tokens. This is also where you will join other domains.",
			Color: 255,
		}, &MenuInfo{
			MenuID: MenuID,
			Author: msg.Author.ID,
			Info: []string{},
		},
	)

	Domains := MainMenu.NewChildMenu(discordgo.MessageEmbed{
		Title: "Domain Management",
		Description: "This is where you set your default domain, add new ones and join other public ones.",
		Color: 255,
	}, MenuButton{
		Description: "This is where you set your default domain, manage your domains, and add new ones.",
		Name: "Domain",
		Emoji: "🕸",
	})
	Domains.AddBackButton()
	Domains.Reactions[MenuButton{
		Emoji: "🗺",
		Name: "Add Domain",
		Description: "This will guide you through adding a domain.",
	}] = AddDomain

	return &MainMenu
}
