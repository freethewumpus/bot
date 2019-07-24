package main

import (
	"github.com/bwmarrin/discordgo"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

func ChangePattern(ChannelID string, MessageID string, menu *EmbedMenu, client *discordgo.Session) {
	defer menu.Display(ChannelID, MessageID, client)

	UserInfo := GetUser(menu.MenuInfo.Author)

	embed := &discordgo.MessageEmbed{
		Title: "Waiting for your pattern...",
		Description: "Please enter your new pattern. It can contain the following characters:```e - Emoji\nn - Number\nc - a-z character```It must  have more than 4 characters but under 150 characters.",
	}
	_, err := client.ChannelMessageEditComplex(&discordgo.MessageEdit{
		ID: MessageID,
		Channel: ChannelID,
		Embed: embed,
	})
	if err != nil {
		return
	}

	UserMessage := WaitForMessage(ChannelID, menu.MenuInfo.Author, 5)
	if UserMessage == nil {
		return
	}

	if len(UserMessage.Content) < 4 || len(UserMessage.Content) > 150 {
		return
	}

	for _, v := range UserMessage.Content {
		switch v {
		case 'e':
		case 'n':
		case 'c': {
			continue
		}
		default:
			return
		}
	}
	UserInfo.NamingScheme = UserMessage.Content

	Update := make(map[string]interface{})
	Update["naming_scheme"] = UserInfo.NamingScheme
	menu.Embed.Description = GeneratePatternDescription(UserInfo)

	err = r.Table("users").Get(UserInfo.Id).Update(&Update).Exec(RethinkConnection)
	if err != nil {
		panic(err)
	}
	InvalidateUserCache(UserInfo.Tokens)
}
