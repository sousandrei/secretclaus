package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type User struct {
	username string
	id       int64
}

var users = make([]User, 0)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func removeDuplicateValues(slice []User) []User {
	keys := make(map[string]bool)
	list := []User{}

	for _, entry := range slice {
		if _, value := keys[entry.username]; !value {
			keys[entry.username] = true
			list = append(list, entry)
		}
	}

	return list
}

func main() {
	bot, err := tg.NewBotAPI(os.Getenv("BOT_TOKEN"))
	checkErr(err)

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tg.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		mid := update.Message.MessageID
		cid := update.Message.Chat.ID
		from := update.Message.From.UserName
		text := update.Message.Text
		mtype := update.Message.Chat.Type

		switch text {
		case "/join":
			join(bot, cid, mid, from, mtype)
		case "/leave":
			leave(bot, cid, mid, from, mtype)
		case "/roda":
			roda(bot, cid, mid, from)
		case "/list":
			list(bot, cid, mid, from)
		default:
		}

	}
}

func leave(bot *tg.BotAPI, cid int64, mid int, from string, mtype string) {
	if mtype != "private" {
		msg := tg.NewMessage(cid, "Manda isso pro bot, nao no grupo")
		msg.ReplyToMessageID = mid

		_, err := bot.Send(msg)
		checkErr(err)

		return
	}

	i := 0
	for range users {
		if users[i].username == from {
			break
		}
		i++
	}

	users = append(users[:i], users[i+1:]...)

	txt := fmt.Sprintf("%s te tirei da brincadeira %d", from, i)

	msg := tg.NewMessage(cid, txt)
	msg.ReplyToMessageID = mid

	_, err := bot.Send(msg)
	checkErr(err)

}

func join(bot *tg.BotAPI, cid int64, mid int, from string, mtype string) {
	if mtype != "private" {
		msg := tg.NewMessage(cid, "Manda isso pro bot, nao no grupo")
		msg.ReplyToMessageID = mid

		_, err := bot.Send(msg)
		checkErr(err)

		return
	}

	user := User{from, cid}
	users = append(users, user)
	users = removeDuplicateValues(users)

	txt := fmt.Sprintf("%s te botei na brincadeira", from)

	msg := tg.NewMessage(cid, txt)
	msg.ReplyToMessageID = mid

	_, err := bot.Send(msg)
	checkErr(err)

}

func roda(bot *tg.BotAPI, cid int64, mid int, from string) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(users), func(i, j int) { users[i], users[j] = users[j], users[i] })

	size := len(users)
	for i := range users {
		var txt string

		if i == size-1 {
			txt = fmt.Sprintf("seu par eh: @%s", users[0].username)
		} else {
			txt = fmt.Sprintf("seu par eh: @%s", users[i+1].username)
		}

		msg := tg.NewMessage(users[i].id, txt)

		_, err := bot.Send(msg)
		checkErr(err)

	}

	txt := "RODANO"

	msg := tg.NewMessage(cid, txt)
	msg.ReplyToMessageID = mid

	_, err := bot.Send(msg)
	checkErr(err)
}

func list(bot *tg.BotAPI, cid int64, mid int, from string) {
	txt := "Ta vazio"
	if len(users) > 0 {
		list := make([]string, 0)

		for i := range users {
			list = append(list, fmt.Sprint("@", users[i].username))
		}

		txt = strings.Join(list, "\n")
	}

	msg := tg.NewMessage(cid, txt)
	msg.ReplyToMessageID = mid

	_, err := bot.Send(msg)
	checkErr(err)
}
