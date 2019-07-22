package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type DomainValidationResponse struct {
	success bool
	message *string
}

func ValidateDomain(Domain string) *string {
	if GetDomain(Domain) != nil {
		err := "The domain already exists in the database."
		return &err
	}
	URL := fmt.Sprintf("http://%s/?domain=%s", os.Getenv("DOMAIN_MANAGER_HOSTNAME"), url.QueryEscape(Domain))
	resp, err := http.Get(URL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		return nil
	} else {
		var ResponseBuffer []byte
		_, err := resp.Body.Read(ResponseBuffer)
		if err != nil {
			panic(err)
		}
		var JSONResponse DomainValidationResponse
		err = json.Unmarshal(ResponseBuffer, JSONResponse)
		if err != nil {
			panic(err)
		}
		return JSONResponse.message
	}
}

func AddDomain(ChannelID string, MessageID string, menu *EmbedMenu, client *discordgo.Session) {
	_ = client.MessageReactionsRemoveAll(ChannelID, MessageID)

	embed := &discordgo.MessageEmbed{
		Title: "Waiting for your domain...",
		Description: "You have 5 minutes to enter your domain in this channel. Simply write it in the format of `example.com`. Make sure before you enter it, you set the A record in your DNS to `" + os.Getenv("CLUSTER_IP") + "`. **Your next message in this channel will be automatically scanned/treated as your domain.**",
	}
	msg, err := client.ChannelMessageEditComplex(&discordgo.MessageEdit{
		ID: MessageID,
		Channel: ChannelID,
		Embed: embed,
	})
	if err != nil {
		return
	}

	UserMessage := WaitForMessage(ChannelID, menu.MenuInfo.Author, 5)
	if UserMessage != nil {
		_ = client.ChannelMessageDelete(ChannelID, UserMessage.ID)
		content := strings.ToLower(UserMessage.Content)
		err := ValidateDomain(content)
		if err != nil {
			embed := &discordgo.MessageEmbed{
				Color: 16711680,
				Title: "Invalid domain",
				Description: fmt.Sprintf("There was an error processing your domain: ```%s```Returning to the domain management page.", err),
			}
			_, err := client.ChannelMessageEditComplex(&discordgo.MessageEdit{
				ID: MessageID,
				Channel: ChannelID,
				Embed: embed,
			})
			if err != nil {
				return
			}
		} else {
			user := GetUser(menu.MenuInfo.Author)
			user.NewDomain(content)
			embed := &discordgo.MessageEmbed{
				Color: 32768,
				Title: "Domain added",
				Description: "Your domain is configured. Returning to the domain management page.",
			}
			_, err := client.ChannelMessageEditComplex(&discordgo.MessageEdit{
				ID: MessageID,
				Channel: ChannelID,
				Embed: embed,
			})
			if err != nil {
				return
			}
		}
	}

	menu.Display(ChannelID, msg.ID, client)
}
