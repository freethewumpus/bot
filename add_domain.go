package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type DomainValidationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func ValidateDomain(Domain string) *string {
	if GetDomain(Domain) != nil {
		err := "The domain already exists in the database."
		return &err
	}

	HttpClient := http.Client{
		Timeout: time.Second * 10,
	}

	URL := fmt.Sprintf("http://%s/?domain=%s", os.Getenv("DOMAIN_MANAGER_HOSTNAME"), url.QueryEscape(Domain))
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		panic(err)
	}
	resp, err := HttpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	JSONResponse := DomainValidationResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &JSONResponse)
	if err != nil {
		panic(err)
	}
	if JSONResponse.Success {
		return nil
	}
	return &JSONResponse.Message
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
				Description: fmt.Sprintf("There was an error processing your domain: ```%s```Returning to the domain management page.", *err),
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

	time.Sleep(time.Duration(5) * time.Second)
	menu.Display(ChannelID, msg.ID, client)
}
