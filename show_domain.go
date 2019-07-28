package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
	"strings"
)

func ShowDomain(domain string) func(ChannelID string, MessageID string, menu *EmbedMenu, client *discordgo.Session) {
	 return func(ChannelID string, MessageID string, OuterMenu *EmbedMenu, client *discordgo.Session) {
		DomainInfo := GetDomain(domain)
		if DomainInfo == nil {
			return
		}
		user := GetUser(OuterMenu.MenuInfo.Author)

		var BucketWarning string
		if DomainInfo.Bucket != nil {
			BucketWarning = "**WARNING: This domain uses a bucket by the domain owner. Freethewump.us cannot be responsible for how the data is handled.**\n"
		}
		var Public string
		var BlackWhiteList string

		if DomainInfo.Public {
			Public = "Yes"
			var UsersNeatened []string
			for _, v := range DomainInfo.Blacklist {
				UsersNeatened = append(UsersNeatened, fmt.Sprintf("<@%s>", v))
			}
			Blacklisted := strings.Join(UsersNeatened, ", ")
			if Blacklisted == "" {
				Blacklisted = "Nobody"
			}
			BlackWhiteList = "**Blacklisted:** " + Blacklisted
		} else {
			Public = "No"
			var UsersNeatened []string
			for _, v := range DomainInfo.Whitelist {
				UsersNeatened = append(UsersNeatened, fmt.Sprintf("<@%s>", v))
			}
			Whitelisted := strings.Join(UsersNeatened, ", ")
			if Whitelisted == "" {
				Whitelisted = "Nobody"
			}
			BlackWhiteList = "**Whitelisted:** " + Whitelisted
		}

		DefaultDomain := "**Default Domain:** "
		if user.Domain == DomainInfo.Id {
			DefaultDomain += "Yes"
		} else {
			DefaultDomain += "No"
		}
		Description := fmt.Sprintf("%s**Owner:** <@%s>\n**Public:** %s\n%s\n%s", BucketWarning, DomainInfo.Owner, Public, BlackWhiteList, DefaultDomain)

		embed := discordgo.MessageEmbed{
			Title: domain,
			Description: Description,
		}

		NewEmbedMenu := NewEmbedMenu(embed, OuterMenu.MenuInfo)
		NewEmbedMenu.parent = OuterMenu
		NewEmbedMenu.Reactions[MenuButton{
			Description: "Goes back a page.",
			Name: "Back",
			Emoji: "⬆",
		}] = func(ChannelID string, MessageID string, _ *EmbedMenu, client *discordgo.Session) {
			OuterMenu.Display(ChannelID, MessageID, client)
		}

		if (!DomainInfo.Public && StringInSlice(OuterMenu.MenuInfo.Author, DomainInfo.Whitelist)) || (
			DomainInfo.Public && !StringInSlice(OuterMenu.MenuInfo.Author, DomainInfo.Blacklist)) {
				if user.Domain != DomainInfo.Id {
					NewEmbedMenu.Reactions[MenuButton{
						Emoji: "🚀",
						Description: "This sets this domain as your default.",
						Name: "Set As Default Domain",
					}] = func(_ string, _ string, _ *EmbedMenu, _ *discordgo.Session) {
						defer ShowDomain(domain)(ChannelID, MessageID, OuterMenu, client)
						NewUser := GetUser(OuterMenu.MenuInfo.Author)
						Update := make(map[string]interface{})
						Update["domain"] = domain
						err := r.Table("users").Get(NewUser.Id).Update(&Update).Exec(RethinkConnection)
						if err != nil {
							panic(err)
						}
						InvalidateUserCache(NewUser.Tokens)
					}
				}
		}

		if DomainInfo.Owner == OuterMenu.MenuInfo.Author {
			if DomainInfo.Public {
				NewEmbedMenu.Reactions[MenuButton{
					Emoji: "🔐",
					Description: "This makes the domain private. For users to use the domain, you will need to whitelist them.",
					Name: "Privatise Domain",
				}] = func(_ string, _ string, _ *EmbedMenu, _ *discordgo.Session) {
					defer ShowDomain(domain)(ChannelID, MessageID, OuterMenu, client)
					Update := make(map[string]interface{})
					Update["public"] = false
					err := r.Table("domains").Get(DomainInfo.Id).Update(&Update).Exec(RethinkConnection)
					if err != nil {
						panic(err)
					}
					InvalidateDomainCache(domain)
				}
				NewEmbedMenu.Reactions[MenuButton{
					Emoji: "👤",
					Description: "This will allow you to blacklist users.",
					Name: "Blacklist Users",
				}] = HandleBlackWhitelist(OuterMenu, false, domain)
			} else {
				NewEmbedMenu.Reactions[MenuButton{
					Emoji: "🔓",
					Description: "This makes the domain public. To stop people using the domain, you will need to blacklist them.",
					Name: "Publicise Domain",
				}] = func(_ string, _ string, _ *EmbedMenu, _ *discordgo.Session) {
					defer ShowDomain(domain)(ChannelID, MessageID, OuterMenu, client)
					Update := make(map[string]interface{})
					Update["public"] = true
					err := r.Table("domains").Get(DomainInfo.Id).Update(&Update).Exec(RethinkConnection)
					if err != nil {
						panic(err)
					}
					InvalidateDomainCache(domain)
				}
				NewEmbedMenu.Reactions[MenuButton{
					Emoji: "👱",
					Description: "This will allow you to whitelist users.",
					Name: "Whitelist Users",
				}] = HandleBlackWhitelist(OuterMenu, true, domain)
			}

			NewEmbedMenu.Reactions[MenuButton{
				Emoji: "🔧",
				Description: "This will allow you to setup BYOB (Bring Your Own Bucket). You need DM's on for this to work.",
				Name: "Bring Your Own Bucket",
			}] = HandleBYOB(domain)
		}

		_ = client.MessageReactionsRemoveAll(ChannelID, MessageID)
		NewEmbedMenu.Display(ChannelID, MessageID, client)
	 }
}
