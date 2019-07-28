package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
)

var RedisConnection *redis.Client
var RethinkConnection *r.Session
var BYOBTemplate []byte

func main() {
	template, err := ioutil.ReadFile("./byob_template.yaml")
	if err != nil {
		panic(err)
	}
	BYOBTemplate = template

	UtilsInit()

	MenuCache = map[string]*EmbedMenu{}

	RedisHost := os.Getenv("REDIS_HOST")
	if RedisHost == "" {
		RedisHost = "localhost:6379"
	}
	RedisPassword := os.Getenv("REDIS_PASSWORD")
	RedisConnection = redis.NewClient(&redis.Options{
		Addr: RedisHost,
		Password: RedisPassword,
		DB: 0,
	})

	RethinkHost := os.Getenv("RETHINK_HOST")
	if RethinkHost == "" {
		RethinkHost = "127.0.0.1:28015"
	}
	RethinkPass := os.Getenv("RETHINK_PASSWORD")
	RethinkUser := os.Getenv("RETHINK_USER")
	if RethinkUser == "" {
		RethinkUser = "admin"
	}
	s, err := r.Connect(r.ConnectOpts{
		Address: RethinkHost,
		Password: RethinkPass,
		Username: RethinkUser,
		Database: "freethewumpus",
	})
	if err != nil {
		panic(err)
	}
	RethinkConnection = s

	discord, err := discordgo.New("Bot " + os.Getenv("TOKEN"))
	if err != nil {
		panic(err)
	}

	DiscordInit(discord)

	err = discord.Open()
	if err != nil {
		panic(err)
	}

	fmt.Println("Bot is running. Always listening.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	_ = discord.Close()
}
