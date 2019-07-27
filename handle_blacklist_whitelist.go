package main

import (
	"github.com/bwmarrin/discordgo"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
	"strings"
)

func IsNumber(Character int32) bool {
	Numbers := "0123456789"
	for _, v := range Numbers {
		if v == Character {
			return true
		}
	}
	return false
}

func HandleBlackWhitelist(OuterMenu *EmbedMenu, Whitelist bool, domain string) func(ChannelID string, MessageID string, _ *EmbedMenu, client *discordgo.Session) {
	return func(ChannelID string, MessageID string, _ *EmbedMenu, client *discordgo.Session) {
		_ = client.MessageReactionsRemoveAll(ChannelID, MessageID)

		defer ShowDomain(domain)(ChannelID, MessageID, OuterMenu, client)

		var ListName string
		if Whitelist {
			ListName = "Whitelist"
		} else {
			ListName = "Blacklist"
		}
		ListNameLower := strings.ToLower(ListName)
		embed := discordgo.MessageEmbed{
			Title: "Waiting for users to " + ListNameLower + "...",
			Description: "Please enter their mentions/IDs with spaces in between.",
		}
		_, err := client.ChannelMessageEditComplex(&discordgo.MessageEdit{
			Embed: &embed,
			ID: MessageID,
			Channel: ChannelID,
		})
		if err != nil {
			return
		}

		msg := WaitForMessage(ChannelID, OuterMenu.MenuInfo.Author, 5)
		if msg == nil {
			return
		}
		_ = client.ChannelMessageDelete(ChannelID, msg.ID)

		split := strings.Split(msg.Content, " ")
		var ValidUsers []*discordgo.User
		for _, v := range split {
			var NumbersOnly string
			for _, x := range v {
				if IsNumber(x) {
					NumbersOnly += string(x)
				}
			}
			DiscordUser, err := client.User(NumbersOnly)
			if err != nil {
				continue
			}
			ValidUsers = append(ValidUsers, DiscordUser)
		}

		DomainInfo := GetDomain(domain)
		if Whitelist {
			WhitelistUsers := DomainInfo.Whitelist
			for _, v := range ValidUsers {
				WhitelistUsers = append(WhitelistUsers, v.ID)
			}
			Update := make(map[string]interface{})
			Update["whitelist"] = Whitelist
			err := r.Table("domains").Get(DomainInfo.Id).Update(&Update).Exec(RethinkConnection)
			if err != nil {
				panic(err)
			}
		} else {
			Blacklist := DomainInfo.Blacklist
			for _, v := range ValidUsers {
				Blacklist = append(Blacklist, v.ID)
			}
			Update := make(map[string]interface{})
			Update["blacklist"] = Blacklist
			err := r.Table("domains").Get(DomainInfo.Id).Update(&Update).Exec(RethinkConnection)
			if err != nil {
				panic(err)
			}
		}

		InvalidateDomainCache(domain)
	}
}
