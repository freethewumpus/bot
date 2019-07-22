package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
	"time"
)

func GenerateToken(user User) string {
	id := uuid.Must(uuid.NewRandom()).String()
	Update := make(map[string]interface{})
	Update["tokens"] = append(user.Tokens, id)
	err := r.Table("users").Get(user.Id).Update(&Update).Exec(RethinkConnection)
	if err != nil {
		panic(err)
	}
	err = r.Table("tokens").Insert(&Token{
		Id: id,
		Uid: user.Id,
	}).Exec(RethinkConnection)
	if err != nil {
		panic(err)
	}
	return id
}

func TokenGenerationEmbed(ChannelID string, MessageID string, menu *EmbedMenu, client *discordgo.Session) {
	_ = client.MessageReactionsRemoveAll(ChannelID, MessageID)
	user := GetUser(menu.MenuInfo.Author)
	token := GenerateToken(user)
	embed := &discordgo.MessageEmbed{
		Color: 32768,
		Title: "Token generated",
		Description: "A token has been generated and DM'd to you. Returning to the token management page.",
	}
	ErrEmbed := &discordgo.MessageEmbed{
		Color: 16711680,
		Title: "DM error",
		Description: "The bot could not DM you. Do you have DM's off or have you blocked the bot?",
	}
	UserChannel, err := client.UserChannelCreate(menu.MenuInfo.Author)
	if err != nil {
		embed = ErrEmbed
	} else {
		_, err = client.ChannelMessageSend(UserChannel.ID, fmt.Sprintf("Your token is `%s`.", token))
	}

	if err != nil {
		embed = ErrEmbed
	}
	_, err = client.ChannelMessageEditComplex(&discordgo.MessageEdit{
		ID: MessageID,
		Channel: ChannelID,
		Embed: embed,
	})
	if err != nil {
		return
	}
	time.Sleep(time.Second * time.Duration(5))
	menu.Display(ChannelID, MessageID, client)
}

func InvalidateAllTokens(user User) {
	InvalidateUserCache(user.Tokens)
	for _, v := range user.Tokens {
		err := r.Table("tokens").Get(v).Delete().Exec(RethinkConnection)
		if err != nil {
			panic(err)
		}
	}
	Update := make(map[string]interface{})
	Update["tokens"] = make([]string, 0)
	err := r.Table("users").Get(user.Id).Update(&Update).Exec(RethinkConnection)
	if err != nil {
		panic(err)
	}
}

func TokenInvalidationEmbed(ChannelID string, MessageID string, menu *EmbedMenu, client *discordgo.Session) {
	_ = client.MessageReactionsRemoveAll(ChannelID, MessageID)
	user := GetUser(menu.MenuInfo.Author)
	InvalidateAllTokens(user)
	embed := &discordgo.MessageEmbed{
		Color: 32768,
		Title: "Tokens invalidated",
		Description: "Your tokens have been invalidated. Returning to the token management page.",
	}
	_, err := client.ChannelMessageEditComplex(&discordgo.MessageEdit{
		ID: MessageID,
		Channel: ChannelID,
		Embed: embed,
	})
	if err != nil {
		return
	}
	time.Sleep(time.Second * time.Duration(5))
	menu.Display(ChannelID, MessageID, client)
}
