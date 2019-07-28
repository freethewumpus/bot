package main

import (
	"bytes"
	"github.com/bwmarrin/discordgo"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
)

type BYOBResult struct {
	S3Region string `yaml:"s3_region"`
	S3Endpoint string `yaml:"s3_endpoint"`
	S3Bucket string `yaml:"s3_bucket"`
	SecretAccessKey string `yaml:"secret_access_key"`
	AccessKeyId string `yaml:"access_key_id"`
}

func HandleBYOB(domain string) func(ChannelID string, MessageID string, menu *EmbedMenu, client *discordgo.Session) {
	return func(ChannelID string, MessageID string, menu *EmbedMenu, client *discordgo.Session) {
		defer ShowDomain(domain)(ChannelID, MessageID, menu.parent, client)
		_ = client.MessageReactionsRemoveAll(ChannelID, MessageID)

		embed := &discordgo.MessageEmbed{
			Title:       "Waiting for your configuration...",
			Description: "A configuration template has been DM'd to you. Please fill it out in the next 15 minutes and DM the file back to the bot.",
		}
		ErrEmbed := &discordgo.MessageEmbed{
			Color:       16711680,
			Title:       "DM error",
			Description: "The bot could not DM you. Do you have DM's off or have you blocked the bot?",
		}

		UserChannel, err := client.UserChannelCreate(menu.MenuInfo.Author)
		if err == nil {
			_, err = client.ChannelMessageSendComplex(UserChannel.ID, &discordgo.MessageSend{
				File: &discordgo.File{
					Name:        "config_template.yaml",
					ContentType: "text/yaml",
					Reader:      bytes.NewReader(BYOBTemplate),
				},
			})
		}

		if err != nil {
			embed = ErrEmbed
		}

		_, err = client.ChannelMessageEditComplex(&discordgo.MessageEdit{
			Embed:   embed,
			ID:      MessageID,
			Channel: ChannelID,
		})

		if err != nil {
			return
		}

		res := WaitForMessage(UserChannel.ID, menu.MenuInfo.Author, 15)
		if res != nil && res.Attachments != nil && len(res.Attachments) == 1 {
			Attachment := res.Attachments[0]
			if 1000000 > Attachment.Size {
				URL := Attachment.URL
				resp, err := http.Get(URL)
				if err != nil {
					panic(err)
				}
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					panic(err)
				}
				var Result BYOBResult
				err = yaml.Unmarshal(body, &Result)
				if err != nil {
					return
				}
				if Result.AccessKeyId == "" || Result.S3Bucket == "" || Result.S3Endpoint == "" || Result.S3Region == "" || Result.SecretAccessKey == "" {
					return
				}
				S3BucketItem := S3Bucket{
					Endpoint:        Result.S3Endpoint,
					Bucket:          Result.S3Bucket,
					AccessKeyId:     Result.AccessKeyId,
					SecretAccessKey: Result.SecretAccessKey,
					Region:          Result.S3Region,
				}
				Update := make(map[string]interface{})
				Update["bucket"] = S3BucketItem
				InvalidateDomainCache(domain)
				err = r.Table("domains").Get(domain).Update(&Update).Exec(RethinkConnection)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}
