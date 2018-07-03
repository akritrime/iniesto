package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	sg "os/signal"
	sc "syscall"

	dg "github.com/bwmarrin/discordgo"
	_ "github.com/lib/pq"
	"github.com/pariz/gountries"
)

var (
	botID       string
	db          *sql.DB
	countries   = gountries.New()
	b           *url.URL
	f           Client
	customFlags = map[string]int{
		"flag_eng": 463313404254486529,
	}
)

const DATE_LAYOUT = "2006-01-02"

func main() {

	b, err := url.Parse("http://api.football-data.org/v2/")
	if err != nil {
		fmt.Println("Problem parsing the base URL:", err)
		return
	}

	ak := os.Getenv("API_KEY")
	if ak == "" {
		fmt.Println("API_KEY NOT SET.")
		return
	}

	f = Client{
		APIKey:  ak,
		BaseURL: b,
		Client:  &http.Client{Timeout: 3 * time.Second},
	}

	fmt.Println("Hello, world!")
	tk := os.Getenv("TOKEN")
	if tk == "" {
		fmt.Println("TOKEN NOT SET.")
		return
	}
	bot, err := dg.New(fmt.Sprintf("Bot %v", tk))

	if err != nil {
		fmt.Println("error in bot")
		return
	}

	user, err := bot.User("@me")
	if err != nil {
		fmt.Println("error in user")
		return
	}

	db, err = getDB()
	if err != nil {
		fmt.Println("error in getting DB", err)
		return
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS guilds (
			id VARCHAR PRIMARY KEY,
			name TEXT,
			prefix TEXT
		)
	`)
	if err != nil {
		fmt.Println("err in creating table guilds", err)
		return
	}

	botID = user.ID

	bot.AddHandler(pingPong)
	bot.AddHandler(hi)
	bot.AddHandler(bye)
	bot.AddHandler(commands)

	err = bot.Open()
	if err != nil {
		fmt.Println("error opening Discord connection,", err)
		return
	}
	defer bot.Close()

	fmt.Println("Bot is running")

	scn := make(chan os.Signal, 1)
	sg.Notify(scn, sc.SIGINT, sc.SIGTERM, os.Interrupt, os.Kill)
	<-scn

}
