package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/rivo/uniseg"
)

const (
	token = ""
	boss  = ""
)

func main() {

	log.Println("piskosender-bot started")

	go listen()
	log.Println("port listening active")

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		err = fmt.Errorf("can't create new bot: %s", err)
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	var ucfg = tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	updates, err := bot.GetUpdatesChan(ucfg)
	if err != nil {
		err = fmt.Errorf("can't get updates channel: %s", err)
		panic(err)
	}
	time.Sleep(time.Millisecond * 500)
	updates.Clear()
	log.Println("piskosender-bot ready")

	banned := make(map[string]bool, 0)
L:
	for {
		select {
		case update := <-updates:
			if !update.Message.IsCommand() {
				continue L
			}
			chatId := update.Message.Chat.ID
			log.Printf("Got new message in chat #%d", chatId)
			var (
				reply    string
				markdown = true
			)
			user := update.Message.From.UserName
			if banned[user] {
				continue L
			}
			switch update.Message.Command() {
			case "start":
				log.Printf("User: '%s' Command: 'start'", user)
				reply = "Ready to get some piskas, sucker?\n/piska for nice, long content"
			case "help":
				log.Printf("User: '%s' Command: 'help'", user)
				reply = "This bot is capable of generating random length male organs\n/piska to try it out"
			case "piska":
				log.Printf("User: '%s' Command: 'piska', reply to: '%v'", user, update.Message.ReplyToMessage)
				if m := update.Message.ReplyToMessage; m != nil {
					reply += "@" + m.From.UserName + ", this one is for you\n"
				}
				reply += "3"
				length := rand.Intn(25)
				if user == boss {
					length = 10 + rand.Intn(40)
				}
				for i := 0; i < length; i++ {
					reply += "="
				}
				reply += ">"
			case "bunner":
				text := getText(update.Message.Text)
				if text == "" {
					reply = nothingToSayBunner
				} else {
					n := longestLine(text)
					offset := 0
					if n%2 == 0 {
						offset = 5
					} else {
						offset = 3
					}
					if n < 5 {
						n = 5
					}
					//верхняя часть баннера
					reply = "|"
					for i := 0; i < n/2+2; i++ {
						reply += "￣"
					}
					reply += "|\n"
					//часть с текстом
					lines := strings.Split(text, "\n")
					for _, line := range lines {
						//отступ от начала
						for i := 0; i < offset; i++ {
							reply += " "
						}
						diff := n - uniseg.GraphemeClusterCount(line)
						log.Printf("DIFF: %d, n=%d, len=%d\n", diff, n, uniseg.GraphemeClusterCount(line))
						//придвинуть строчку к центру баннера
						for i := 0; i < diff/2; i++ {
							reply += "  "
						}
						reply += line + "\n"
					}
					reply += "|"
					//нижняя часть баннера
					for i := 0; i < n/2+2; i++ {
						reply += "＿"
					}
					reply += "|\n"
					//кролик, придвинутый к центру баннера
					for i := 0; i < (n-7)/2; i++ {
						reply += "  "
					}
					reply += `   (\＿/) ||` + "\n"
					for i := 0; i < (n-7)/2; i++ {
						reply += "  "
					}
					reply += `   (•ㅅ•) ||` + "\n"
					for i := 0; i < (n-7)/2; i++ {
						reply += "  "
					}
					reply += `   /   づ` + "\n"
					//markdown = false
				}
			case "kitty":
				if m := update.Message.ReplyToMessage; m != nil {
					user = m.From.UserName
				}
				reply = heartCat + "\n@" + user + ", love you"
				markdown = false
			case "cat":
				reply = hugeCat
			case "wat":
				reply = pokerCat
			case "ban":
				if user != boss {
					continue L
				}
				text := getText(update.Message.Text)
				banned[text] = true
				reply = "@" + text + ", задолбал, забанен"
			case "unban":
				if user != boss {
					continue L
				}
				text := getText(update.Message.Text)
				banned[text] = false
				reply = "@" + text + ", помилован, разбанен"
			default:
				continue L
			}
			log.Println("reply: ", reply)
			msg := tgbotapi.NewMessage(chatId, reply)
			if markdown {
				msg.ParseMode = "markdown"
			}
			if msg.Text == "" {
				continue L
			}
			m, err := bot.Send(msg)
			if err != nil {
				log.Println("Error on message sending: ", err, "\nMessage: ", m)
			}
		}
	}
}

func listen() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("no such env var")
	}
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("can't listen to port ERROR: ", err)
	}
}

func getText(data string) string {
	exp := `\/[^\s]*[@]?[^\s]* *\n*` //match a bot command
	re := regexp.MustCompile(exp)
	newData := re.ReplaceAllString(data, "")
	log.Printf("GOT: '%s'", data)
	log.Printf("RETURNING: '%s'\n", newData)
	return newData
}

func longestLine(data string) int {
	a := strings.Split(data, "\n")
	log.Printf("%d LINES: %+v", len(a), a)
	max := 0
	for _, s := range a {
		if uniseg.GraphemeClusterCount(s) > max {
			max = uniseg.GraphemeClusterCount(s)
		}
	}
	log.Printf("MAXLINE: %d", max)
	return max
}
