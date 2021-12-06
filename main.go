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

		if text == "/list" {
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

			_, err = bot.Send(msg)
			checkErr(err)

			continue
		}

		if text == "/join" {
			if update.Message.Chat.Type != "private" {
				msg := tg.NewMessage(cid, "Manda isso pro bot, nao no grupo")
				msg.ReplyToMessageID = mid

				_, err = bot.Send(msg)
				checkErr(err)

				continue
			}

			user := User{from, cid}
			users = append(users, user)

			txt := fmt.Sprintf("%s te botei na brincadeira", from)

			msg := tg.NewMessage(cid, txt)
			msg.ReplyToMessageID = mid

			_, err = bot.Send(msg)
			checkErr(err)

			continue
		}

		if text == "/roda" {
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

				_, err = bot.Send(msg)
				checkErr(err)

			}

			txt := "RODANO"

			msg := tg.NewMessage(cid, txt)
			msg.ReplyToMessageID = mid

			_, err = bot.Send(msg)
			checkErr(err)

			continue
		}

	}
}
