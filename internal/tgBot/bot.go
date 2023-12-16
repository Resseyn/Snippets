package tgBot

import (
	"SnippetsTESTBYGUIDE/internal/loggers"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"time"
)

var ActiveSessionIsCreatingOrEditing = make(map[int64]bool)  //key = chatID
var UserQueryChannels = make(map[int64]chan tgbotapi.Update) //key = chatID
var ChatCurrStack = make(map[int64]int)                      //key = chatID

func BotStart() {
	bot, err := tgbotapi.NewBotAPI("6467098865:AAHByMBybrT_pFOjySUOg960m6YiW7D7B4Y")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)
	go func() {
		for {
			CheckIfSomeSnippetGoesExpired(bot)
			time.Sleep(31 * time.Minute)
		}
	}()

	for update := range updates {
		if update.CallbackQuery == nil {
			if ActiveSessionIsCreatingOrEditing[update.Message.Chat.ID] {
				if update.Message.Text != "" && !update.Message.IsCommand() && update.CallbackQuery == nil {
					if ActiveSessionIsCreatingOrEditing[update.Message.Chat.ID] {
						UserQueryChannels[update.Message.Chat.ID] <- update
						continue
					}
				} else {
					bot.Send(tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID))
					continue
				}
			}
			switch update.Message.Command() {
			case "start":
				go startCommand(bot, &update)
			case "create":
				go createCommand(bot, update.Message.Chat.ID)
			case "showlatest":
				ChatCurrStack[update.Message.Chat.ID] = 0
				go showLatest(bot, update.Message.Chat.ID, true)
			case "showlatestfast":
				ChatCurrStack[update.Message.Chat.ID] = 0
				go showLatest(bot, update.Message.Chat.ID, false)
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
				msg.Text = "What?"
				bot.Send(msg)
			}
			bot.Send(tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID))
		} else {
			if ActiveSessionIsCreatingOrEditing[update.CallbackQuery.Message.Chat.ID] {
				if ActiveSessionIsCreatingOrEditing[update.CallbackQuery.Message.Chat.ID] {
					UserQueryChannels[update.CallbackQuery.Message.Chat.ID] <- update
					continue
				}
			}
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			switch callback.Text {
			case "shownextlatest":
				ChatCurrStack[update.CallbackQuery.Message.Chat.ID] += 5
				go ReshowList(bot, update.CallbackQuery.Message.Chat.ID)
			case "showprevlatest":
				ChatCurrStack[update.CallbackQuery.Message.Chat.ID] -= 5
				go ReshowList(bot, update.CallbackQuery.Message.Chat.ID)
			default:
				var data JsonWithCommandAndData
				err := json.Unmarshal([]byte(callback.Text), &data)
				if err != nil {
					loggers.Logger.Println(err)
					continue
				}
				switch data.Command {
				case "unboxsnippet":
					go func() {
						err := unboxSnippet(bot, update.CallbackQuery.Message.Chat.ID, &data)
						if err != nil {
							loggers.Logger.Println(err)
						}
					}()
				case "update":
					go func() {
						err := updateSnippet(bot, update.CallbackQuery.Message.Chat.ID, &data)
						if err != nil {
							loggers.Logger.Println(err)
						}
					}()
				case "delete":
					go func() {
						err := deleteSnippet(bot, update.CallbackQuery.Message.Chat.ID, &data)
						if err != nil {
							loggers.Logger.Println(err)
						}
					}()
				case "extend":
					go func() {
						err := extendSnippet(bot, update.CallbackQuery.Message.Chat.ID, &data)
						if err != nil {
							loggers.Logger.Println(err)
						}
					}()
				case "close":
					go func() {
						err := closeSnippet(bot, update.CallbackQuery.Message.Chat.ID, &data)
						if err != nil {
							loggers.Logger.Println(err)
						}
						ReshowList(bot, update.CallbackQuery.Message.Chat.ID)
					}()
				}
			}
		}
	}
}
