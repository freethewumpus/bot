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
		Name: "Domain Management",
		Emoji: "🕸",
	})
	Domains.AddBackButton()
	Domains.Reactions[MenuButton{
		Emoji: "🗺",
		Name: "Add Domain",
		Description: "This will guide you through adding a domain.",
	}] = AddDomain
	OwnedDomains := GetUser(msg.Author.ID).GetOwnedDomains()
	Domains.Reactions[MenuButton{
		Emoji: "👥",
		Name: "Get Owned Domains",
		Description: "This will get all owned domains.",
	}] = DomainPages("Owned Domains", OwnedDomains, 0, ShowDomain)

	Tokens := MainMenu.NewChildMenu(discordgo.MessageEmbed{
		Title: "Tokens Management",
		Description: "This is where you can add or revoke tokens.",
		Color: 255,
	}, MenuButton{
		Description: "This is where you can add or revoke tokens.",
		Name: "Token Management",
		Emoji: "🎟",
	})
	Tokens.AddBackButton()
	Tokens.Reactions[MenuButton{
		Emoji: "🥊",
		Name: "Revoke Tokens",
		Description: "This will revoke all your tokens.",
	}] = TokenInvalidationEmbed
	Tokens.Reactions[MenuButton{
		Emoji: "🎟",
		Name: "Generate Token",
		Description: "This will generate a token. This requires DM's to be on.",
	}] = TokenGenerationEmbed

	UserInfo := GetUser(msg.Author.ID)

	FileNamingPattern := MainMenu.NewChildMenu(discordgo.MessageEmbed{
		Title: "File Naming Pattern",
		Description: GeneratePatternDescription(UserInfo),
		Color: 255,
	}, MenuButton{
		Description: "This is where you can set your file naming pattern.",
		Name: "File Naming Pattern",
		Emoji: "🗄",
	})
	FileNamingPattern.AddBackButton()
	FileNamingPattern.Reactions[MenuButton{
		Description: "This will allow you to change your file naming pattern.",
		Name: "Change Pattern",
		Emoji: "🗄",
	}] = ChangePattern

	return &MainMenu
}
