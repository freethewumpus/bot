package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func DomainPages(
	title string, domains []Domain, after int,
	function func(item string) func(ChannelID string, MessageID string, menu *EmbedMenu, client *discordgo.Session)) func(ChannelID string, MessageID string, menu *EmbedMenu, client *discordgo.Session) {
		return func(ChannelID string, MessageID string, menu *EmbedMenu, client *discordgo.Session) {
			_ = client.MessageReactionsRemoveAll(ChannelID, MessageID)

			embed := discordgo.MessageEmbed{
				Title: title,
			}
			NewEmbedMenu := NewEmbedMenu(embed, menu.MenuInfo)
			NewEmbedMenu.parent = menu
			NewEmbedMenu.AddBackButton()

			var DomainPart []Domain
			total := 0
			for i, v := range domains {
				if i + 1 > after {
					DomainPart = append(DomainPart, v)
					total++
				}
			}

			if len(domains) > 5 && after != 0 {
				NewEmbedMenu.Reactions[MenuButton{
					Name: "Back Page",
					Description: "Go back a page.",
					Emoji: "◀",
				}] = DomainPages(title, domains, after + 5, function)
			}

			if len(domains) > 5 && len(DomainPart) == 5 {
				NewEmbedMenu.Reactions[MenuButton{
					Name: "Forward Page",
					Description: "Go forward a page.",
					Emoji: "▶",
				}] = DomainPages(title, domains, after + 5, function)
			}

			letters := []string{"🇦", "🇧", "🇨", "🇩", "🇪"}
			for i, v := range DomainPart {
				NewEmbedMenu.Reactions[MenuButton{
					Name: v.Id,
					Description: fmt.Sprintf("Owned by <@%s>.", v.Owner),
					Emoji: letters[i],
				}] = function(v.Id)
			}

			NewEmbedMenu.Display(ChannelID, MessageID, client)
		}
}
