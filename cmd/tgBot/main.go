package main

import (
	// "fmt"
	"TODobot/pkg/commands"
	"TODobot/pkg/task"
	"TODobot/pkg/users"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/skinass/telegram-bot-api/v5"
)

// https://api.telegram.org/bot5476117204:AAGwDhxItW4ieg6_xrOFKp-RVCVYgk5Po64/getUpdates

var  (
	BotToken = "5476117204:AAGwDhxItW4ieg6_xrOFKp-RVCVYgk5Po64"

	WebHookURL = "https://ac35-85-143-112-90.ngrok-free.app"
)





func main() {
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Fatalf("NewBotAPI failed: %s", err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s\n", bot.Self.UserName)

	wh, err := tgbotapi.NewWebhook(WebHookURL)
	if err != nil {
		log.Fatalf("NewWebhook failed: %s", err)
	}

	_, err = bot.Request(wh)
	if err != nil {
		log.Fatalf("SetWebhook failed: %s", err)
	}

	updates := bot.ListenForWebhook("/")

	http.HandleFunc("/state", func(w http.ResponseWriter, r *http.Request) {
		_, err1 := w.Write([]byte("all is working"))
		if err1 != nil {
			log.Println("не работает")
			return
		}
	})
	// var i int = 10
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	userRepo := users.NewMemoryRepo()
	taskRepo := task.NewMemoryRepo()

	go func() {
		log.Fatalln("http err:", http.ListenAndServe(":"+port, nil))
	}()
	log.Println("start listen :" + port)

	
	for update := range updates {
		if update.Message == nil {
			continue
		}
		log.Println("message ", update.Message.Text)
		


		currUser, err := userRepo.GetUser(update.Message.From.UserName)
		if err != nil {
			currUser = users.User{UserName: update.Message.From.UserName, ChatId: update.Message.Chat.ID}
			userRepo.AddNewUser(currUser)
		}

		if update.Message.IsCommand() {
			commands.ForCommand(*bot, currUser, update, taskRepo, userRepo)
		} else {
			_, err1 := bot.Send(tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"Привет, напиши /help для команд",
			))
			if err1 != nil {
				log.Fatalf("SetWebhook failed: %s", err1)
			}
		}

	}
}
