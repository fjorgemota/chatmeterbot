package main

import "sync"
import "gopkg.in/telegram-bot-api.v4"

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
