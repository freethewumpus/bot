package main

import (
	"github.com/bwmarrin/discordgo"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

func CreateNewMenu(MenuID string, msg discordgo.Message) *EmbedMenu {
	MainMenu := NewEmbedMenu(
		discordgo.MessageEmbed{
			Title:       "Freethewump.us Manager",
			Description: "Using this bot, you can manage your domains and tokens. This is also where you will join other domains.",
			Color:       255,
		}, &MenuInfo{
			MenuID: MenuID,
			Author: msg.Author.ID,
			Info:   []string{},
		},
	)

	user := GetUser(msg.Author.ID)

	EncryptionName := "Enable Encryption"
	if user.Encryption {
		EncryptionName = "Disable Encryption"
	}
	MainMenu.Reactions.Add(MenuReaction{
		button: MenuButton{
			Description: "This will allow you to toggle encryption.",
			Name:        EncryptionName,
			Emoji:       "🔑",
		},
		function: func(ChannelID string, MessageID string, menu *EmbedMenu, client *discordgo.Session) {
			_ = client.MessageReactionsRemoveAll(ChannelID, MessageID)
			defer CreateNewMenu(MenuID, msg).Display(ChannelID, MessageID, client)

			user.Encryption = !user.Encryption
			r.Table("users").Get(user.Id).Update(&map[string]interface{}{
				"encryption": user.Encryption,
			}).Exec(RethinkConnection)
		},
	})

	Domains := MainMenu.NewChildMenu(discordgo.MessageEmbed{
		Title:       "Domain Management",
		Description: "This is where you set your default domain, add new ones and join other public ones.",
		Color:       255,
	}, MenuButton{
		Description: "This is where you set your default domain, manage your domains, and add new ones.",
		Name:        "Domain Management",
		Emoji:       "🕸",
	})
	Domains.AddBackButton()
	Domains.Reactions.Add(MenuReaction{
		button: MenuButton{
			Emoji:       "🗺",
			Name:        "Add Domain",
			Description: "This will guide you through adding a domain.",
		},
		function: AddDomain,
	})
	Domains.Reactions.Add(MenuReaction{
		button: MenuButton{
			Emoji:       "👥",
			Name:        "Get Owned Domains",
			Description: "This will get all owned domains.",
		},
		function: DomainPages("Owned Domains", user.GetOwnedDomains, 0, ShowDomain),
	})
	Domains.Reactions.Add(MenuReaction{
		button: MenuButton{
			Emoji:       "👱",
			Name:        "Get Whitelisted Domains",
			Description: "This will get all whitelisted domains.",
		},
		function: DomainPages("Whitelisted Domains", user.GetWhitelistedDomains, 0, ShowDomain),
	})
	Domains.Reactions.Add(MenuReaction{
		button: MenuButton{
			Emoji:       "🌐",
			Name:        "Get Public Domains",
			Description: "This will get all public domains.",
		},
		function: DomainPages("Public Domains", GetPublicDomains, 0, ShowDomain),
	})

	Tokens := MainMenu.NewChildMenu(discordgo.MessageEmbed{
		Title:       "Tokens Management",
		Description: "This is where you can add or revoke tokens.",
		Color:       255,
	}, MenuButton{
		Description: "This is where you can add or revoke tokens.",
		Name:        "Token Management",
		Emoji:       "🎟",
	})
	Tokens.AddBackButton()
	Tokens.Reactions.Add(MenuReaction{
		button: MenuButton{
			Emoji:       "🥊",
			Name:        "Revoke Tokens",
			Description: "This will revoke all your tokens.",
		},
		function: TokenInvalidationEmbed,
	})
	Tokens.Reactions.Add(MenuReaction{
		button: MenuButton{
			Emoji:       "🎟",
			Name:        "Generate Token",
			Description: "This will generate a token. This requires DM's to be on.",
		},
		function: TokenGenerationEmbed,
	})

	UserInfo := GetUser(msg.Author.ID)

	FileNamingPattern := MainMenu.NewChildMenu(discordgo.MessageEmbed{
		Title:       "File Naming Pattern",
		Description: GeneratePatternDescription(UserInfo),
		Color:       255,
	}, MenuButton{
		Description: "This is where you can set your file naming pattern.",
		Name:        "File Naming Pattern",
		Emoji:       "🗄",
	})
	FileNamingPattern.AddBackButton()
	FileNamingPattern.Reactions.Add(MenuReaction{
		button: MenuButton{
			Description: "This will allow you to change your file naming pattern.",
			Name:        "Change Pattern",
			Emoji:       "🗄",
		},
		function: ChangePattern,
	})

	return &MainMenu
}
