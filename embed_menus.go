package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

var MenuCache map[string]*EmbedMenu

type MenuInfo struct {
	MenuID string
	Author string
	Info []string
}

type MenuButton struct {
	Emoji string
	Name string
	Description string
}

type EmbedMenu struct {
	Reactions map[MenuButton]func(ChannelID string, MessageID string, menu *EmbedMenu, client *discordgo.Session)
	parent *EmbedMenu
	Embed *discordgo.MessageEmbed
	MenuInfo *MenuInfo
}

func (emm EmbedMenu) Display(ChannelID string, MessageID string, client *discordgo.Session) *error {
	MenuCache[emm.MenuInfo.MenuID] = &emm

	EmbedCopy := emm.Embed
	EmbedCopy.Footer = &discordgo.MessageEmbedFooter{
		Text: emm.MenuInfo.MenuID,
	}
	for k := range emm.Reactions {
		EmbedCopy.Fields = append(EmbedCopy.Fields, &discordgo.MessageEmbedField{
			Name: fmt.Sprintf("%s %s", k.Emoji, k.Name),
			Value: k.Description,
			Inline: false,
		})
	}
	_, err := client.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Embed: EmbedCopy,
		ID: MessageID,
		Channel: ChannelID,
	})
	if err != nil {
		return &err
	}
	for k := range emm.Reactions {
		err := client.MessageReactionAdd(ChannelID, MessageID, k.Emoji)
		if err != nil {
			return &err
		}
	}
	return nil
}

func (emm EmbedMenu) NewChildMenu(embed discordgo.MessageEmbed, item MenuButton) *EmbedMenu {
	NewEmbedMenu := NewEmbedMenu(embed, emm.MenuInfo)
	NewEmbedMenu.parent = &emm
	emm.Reactions[item] = func(ChannelID string, MessageID string, _ *EmbedMenu, client *discordgo.Session) {
		_ = client.MessageReactionsRemoveAll(ChannelID, MessageID)
		NewEmbedMenu.Display(ChannelID, MessageID, client)
	}
	return &NewEmbedMenu
}

func (emm EmbedMenu) AddBackButton() {
	emm.Reactions[MenuButton{
		Description: "Goes back a page.",
		Name: "Back",
		Emoji: "⬆",
	}] = func(ChannelID string, MessageID string, _ *EmbedMenu, client *discordgo.Session) {
		_ = client.MessageReactionsRemoveAll(ChannelID, MessageID)
		emm.parent.Display(ChannelID, MessageID, client)
	}
}

func NewEmbedMenu(embed discordgo.MessageEmbed, info *MenuInfo) EmbedMenu {
	menu := EmbedMenu{
		Reactions: map[MenuButton]func(ChannelID string, MessageID string, menu *EmbedMenu, client *discordgo.Session){},
		Embed: &embed,
		MenuInfo: info,
	}
	return menu
}

func HandleMenuReactionEdit(client *discordgo.Session, reaction *discordgo.MessageReactionAdd, MenuID string) {
	_ = client.MessageReactionRemove(reaction.ChannelID, reaction.MessageID, reaction.Emoji.Name, reaction.UserID)
	menu := MenuCache[MenuID]
	if menu == nil {
		return
	}

	if menu.MenuInfo.Author != reaction.UserID {
		return
	}

	for k, v := range menu.Reactions {
		if k.Emoji == reaction.Emoji.Name {
			v(reaction.ChannelID, reaction.MessageID, menu, client)
			return
		}
	}
}
