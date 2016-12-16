package main

import (
	"log"
	"os"
	"sync"
	"time"

	"gopkg.in/telegram-bot-api.v4"
)

func shouldRegisterMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
	return update.Message != nil && update.Message.From.ID != bot.Self.ID && (update.Message.Chat.IsGroup() || update.Message.Chat.IsSuperGroup())
}

func ticker(bot *tgbotapi.BotAPI, stickerID string, interval time.Duration, updates <-chan tgbotapi.Update, w *sync.WaitGroup) error {
	log.Printf("Starting ticker..")
	groups := make(map[int64]time.Time)
	tick := time.Tick(interval)
	run := true
	for run {
		select {
		case update, ok := <-updates:
			run = ok
			if !shouldRegisterMessage(bot, update) {
				continue
			}
			log.Printf("Received message on %d\n", update.Message.Chat.ID)
			groups[update.Message.Chat.ID] = update.Message.Time().Add(interval)
		case <-tick:
			now := time.Now()
			log.Println("Processing tick..")
			for chatID, lastMessage := range groups {
				if lastMessage.After(now) {
					continue
				}
				log.Printf("Sending sticker to %d\n", chatID)
				sticker := tgbotapi.NewStickerShare(chatID, stickerID)
				_, err := bot.Send(sticker)
				if err != nil {
					log.Print(err)
				}
				delete(groups, chatID)
			}
		}
	}
	log.Println("Finishing ticker")
	w.Done()
	return nil
}

func broadcast(w *sync.WaitGroup, source <-chan tgbotapi.Update, dest ...chan tgbotapi.Update) {
	for msg := range source {
		for _, target := range dest {
			target <- msg
		}
	}
	for _, target := range dest {
		close(target)
	}
	w.Done()
}
func main() {
	var updates <-chan tgbotapi.Update
	token := os.Getenv("BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err == nil {
		log.Printf("Authorized on account %s", bot.Self.UserName)
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 3600
		updates, err = bot.GetUpdatesChan(u)
	}
	if err != nil {
		log.Fatal(err)
	}
	updates1 := make(chan tgbotapi.Update, 100)
	updates2 := make(chan tgbotapi.Update, 100)
	waiter := &sync.WaitGroup{}
	waiter.Add(1)
	go broadcast(waiter, updates, updates1, updates2)
	// Silence sticker:
	go ticker(bot, "BQADAQADNwQAAiQ2IAgTxKl54avixQI", time.Minute*10, updates1, waiter)
	// RIP sticker:
	go ticker(bot, "BQADAQADNQQAAiQ2IAjKJvQOmXvP-QI", time.Minute*30, updates2, waiter)
	waiter.Wait()
}
