package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
)

func OnReady(client *discordgo.Session, _ *discordgo.Ready) {
	err := client.UpdateStatus(0, "freethewump.us")
	if err != nil {
		fmt.Println("Error setting game: ", err)
	}
}

func OnMessage(client *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.Author.Bot {
		return
	}

	mentions := msg.Mentions
	if len(mentions) == 1 && mentions[0].ID == client.State.User.ID {
		go func() {
			OpenMenu(client, msg)
		}()
	}
	go func() {
		if msg.Content == "echo" {
			msg := WaitForMessage(msg.ChannelID, msg.Author.ID, 0)
			fmt.Println(msg)
		}
	}()
	MessageWaitHandler(msg.Message)
}

func OpenMenu(client *discordgo.Session, msg *discordgo.MessageCreate) {
	MenuID := uuid.Must(uuid.NewRandom()).String()
	m, err := client.ChannelMessageSendComplex(msg.ChannelID, &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title: "Loading...",
		},
	})
	if err != nil {
		return
	}
	Menu := CreateNewMenu(MenuID, *msg.Message)
	Menu.Display(msg.ChannelID, m.ID, client)
}

func OnReactionAdd(client *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
	message, err := client.ChannelMessage(reaction.ChannelID, reaction.MessageID)
	if err != nil {
		return
	}
	user, err := client.User(reaction.UserID)
	if err != nil {
		return
	}
	if user.Bot {
		return
	}
	if message.Author.ID == client.State.User.ID && len(message.Embeds) == 1 && message.Embeds[0].Footer != nil {
		MenuID := message.Embeds[0].Footer.Text
		go func() {
			HandleMenuReactionEdit(client, reaction, MenuID)
		}()
	}
}
