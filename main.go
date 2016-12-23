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

func ticker(bot *tgbotapi.BotAPI, stickerID string, interval time.Duration, updates <-chan tgbotapi.Update, w *sync.WaitGroup) {
	log.Printf("Starting ticker..")
	defer log.Println("Finishing ticker")
	defer w.Done()
	groups := make(map[int64]time.Time)
	counters := make(map[int64]int)
	tick := time.Tick(interval)
Run:
	for {
		select {
		case update, ok := <-updates:
			if !ok {
				break Run
			}
			if !shouldRegisterMessage(bot, update) {
				continue
			}
			groups[update.Message.Chat.ID] = update.Message.Time().Add(interval)
			counters[update.Message.Chat.ID]++
			log.Printf("Received message on %d (%d messages counted)\n", update.Message.Chat.ID, counters[update.Message.Chat.ID])
		case <-tick:
			now := time.Now()
			log.Println("Processing tick..")
			for chatID, lastMessage := range groups {
				if lastMessage.After(now) || counters[chatID] < int((interval/time.Minute)*2) {
					continue
				}
				log.Printf("Sending sticker to %d\n", chatID)
				sticker := tgbotapi.NewStickerShare(chatID, stickerID)
				_, err := bot.Send(sticker)
				if err != nil {
					log.Print(err)
				}
				delete(groups, chatID)
				delete(counters, chatID)
			}
		}
	}
}

func main() {
	var updates <-chan tgbotapi.Update
	token := os.Getenv("BOT_TOKEN")
	if len(token) == 0 {
		log.Fatal("You need to define a BOT_TOKEN environment variable")
	}
	bot, err := tgbotapi.NewBotAPI(token)
	if err == nil {
		log.Printf("Authorized on account %s", bot.Self.UserName)
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 86400
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
