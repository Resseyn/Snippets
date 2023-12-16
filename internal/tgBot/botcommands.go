package tgBot

import (
	"SnippetsTESTBYGUIDE/internal/database"
	"SnippetsTESTBYGUIDE/internal/loggers"
	"SnippetsTESTBYGUIDE/pkg/models"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
	"time"
)

func startCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	if _, err := database.Users.Insert(update.Message.From.ID, update.Message.Chat.ID); err == nil {
		msg.Text = "Добро пожаловать епта!!!!!!!!!!!!!!!!!!!"
		bot.Send(msg)
		stc := tgbotapi.NewStickerShare(update.Message.Chat.ID, "CAACAgIAAxkBAAEnDwxlNtf3b-ic4psseS-Vget0Ok9RDAACJxcAAqpO0UpJ3hYkRL3VwTAE")
		bot.Send(stc)
	} else {
		if us, _ := database.Users.Get(update.Message.From.ID); us.ChatID != update.Message.Chat.ID {
			err := database.Users.Update(us.UserID, update.Message.Chat.ID)
			if err != nil {
				loggers.Logger.Println(err)
			}
		}
	}
	msg.Text = "Я всегда с тобой)"
	bot.Send(msg)
	stc := tgbotapi.NewStickerShare(update.Message.Chat.ID, "CAACAgIAAxkBAAEm_iBlMtQ1oREk6uZElKOIF0j8zxEtlgACwxgAAmtW0Eoimf1k6tVxYjAE")
	bot.Send(stc)
}
func createCommand(bot *tgbotapi.BotAPI, chatID int64) {
	ActiveSessionIsCreatingOrEditing[chatID] = true
	UserQueryChannels[chatID] = make(chan tgbotapi.Update)
	msg := tgbotapi.NewMessage(chatID, "")
	if CurrentShownLatestListID != 0 {
		_, err := bot.Send(tgbotapi.NewDeleteMessage(chatID, CurrentShownLatestListID))
		if err != nil {
			loggers.Logger.Println(err)
			msg.Text = "Error idi nax"
			bot.Send(msg)
		}
	}
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	msg.Text = "Вы захотели создать записку. Введите заголовок"
	msgs := make([]tgbotapi.Message, 6, 6)
	msgs[0], _ = bot.Send(msg)
	status := make([]string, 0, 3)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	for updat := range UserQueryChannels[chatID] {
		if updat.CallbackQuery == nil && updat.Message.Chat.ID == chatID && !updat.Message.IsCommand() {
			msgs[3+len(status)] = *updat.Message
			status = append(status, updat.Message.Text)
			if len(status) == 1 {
				msg.Text = "Введите Содержание"
				msgs[1], _ = bot.Send(msg)
			} else if len(status) == 2 {
				msg.Text = "Введите время хранения, в часах"
				msgs[2], _ = bot.Send(msg)
			} else if len(status) == cap(status) {
				break
			}
		} else {
			continue
		}
	}
	for i := 0; i < 6; i++ {
		bot.Send(tgbotapi.NewDeleteMessage(chatID, msgs[i].MessageID))
	}
	msg.Text = "Запись создана!"
	var creator models.User
	database.Users.DB.Table("users").Where("chat_id = ?", chatID).First(&creator)
	_, err := database.Snippets.Insert(creator.UserID, status[0], status[1], status[2])
	if err != nil {
		loggers.Logger.Println(err)
		return
	}
	bot.Send(msg)
	stc := tgbotapi.NewStickerShare(chatID, "CAACAgIAAxkBAAEm_k5lMt0x4dcBqFEWV_uVGZ4mkDclbQACuhkAAtaP2Uq75lOvlc29iTAE")
	bot.Send(stc)
	ActiveSessionIsCreatingOrEditing[chatID] = false
	UserQueryChannels[chatID] = nil
}
func showLatest(bot *tgbotapi.BotAPI, chatID int64, isFullyShown bool) {
	bot.Send(tgbotapi.NewDeleteMessage(chatID, CurrentShownLatestListID))
	database.Snippets.DeleteExpired()
	msg := tgbotapi.NewMessage(chatID, "")
	ClearAllShownSnippets(bot, &msg)
	var creator models.User
	database.Users.DB.Model(&models.User{}).Where("chat_id = ?", chatID).Find(&creator)
	var count int64
	database.Snippets.DB.Table("snippets").Where("created_by = ?", creator.UserID).Count(&count)
	if count == 0 {
		msg.Text = "Записей еще нет🤔"
		bot.Send(msg)
		stc := tgbotapi.NewStickerShare(chatID, "CAACAgIAAxkBAAEnFAJlN-2NcfFQZcpEfB-K7tVLEhVTxwACfxgAAuje0EopMHU5JWuFMzAE")
		bot.Send(stc)
	}
	latest, err := database.Snippets.Latest(ChatCurrStack[chatID], 5, creator.UserID)
	if err != nil {
		loggers.Logger.Println(err)
		msg.Text = err.Error()
		bot.Send(msg)
		return
	}
	msg.Text = ""
	timeNow := time.Now()
	if isFullyShown {
		for i, snippet := range latest {
			timeForSnippet := snippet.Expires.Sub(timeNow).Truncate(time.Hour).Hours()
			var contentToShow string
			if len(snippet.Content) >= 100 {
				contentToShow = snippet.Content[0:100] + "..."
			} else {
				contentToShow = snippet.Content
			}
			msg.Text += fmt.Sprintf("%d - <b>%s</b>\n        %s\n                Пропадет через <b>%.0f</b> часов\n", i+1+ChatCurrStack[chatID], snippet.Title, contentToShow, timeForSnippet)
		}
	} else {
		for i, snippet := range latest {
			timeForSnippet := snippet.Expires.Sub(timeNow).Truncate(time.Hour).Hours()
			var contentToShow string
			if len(snippet.Content) >= 100 {
				contentToShow = snippet.Content[0:100] + "..."
			} else {
				contentToShow = snippet.Content
			}
			msg.Text += fmt.Sprintf("%d - <b>%s</b>\n    Осталось <b>%.0f</b> часов\n", i+1+ChatCurrStack[chatID], contentToShow, timeForSnippet)
		}
	}
	toShow := min(count-int64(ChatCurrStack[chatID]), 5)
	dates := make([][]byte, 0, 5)
	for i := 0; int64(i) < toShow; i++ {
		data, err := json.Marshal(JsonWithCommandAndData{
			Command: "unboxsnippet",
			ID:      latest[i].ID,
		})
		if err != nil {
			loggers.Logger.Println(err)
			msg.Text = err.Error()
			bot.Send(msg)
			return
		}
		dates = append(dates, data)
	}
	row := make([]tgbotapi.InlineKeyboardButton, toShow, toShow)
	for i, data := range dates {
		row[i] = tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d", i+ChatCurrStack[chatID]+1), string(data))
	}
	if (toShow < 5 || count-int64(ChatCurrStack[chatID]+5) == 0) && count > 5 {
		SnippetSSSKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			row,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Показать предыдущие", "showprevlatest"),
			),
		)
		msg.ParseMode = "HTML"
		msg.ReplyMarkup = SnippetSSSKeyboard
		sended, _ := bot.Send(msg)
		CurrentShownLatestListID = sended.MessageID
	} else if ChatCurrStack[chatID] >= 5 {
		SnippetSSSKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			row,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Показать предыдущие", "showprevlatest"),
				tgbotapi.NewInlineKeyboardButtonData("Показать еще", "shownextlatest"),
			),
		)
		msg.ParseMode = "HTML"
		msg.ReplyMarkup = SnippetSSSKeyboard
		sended, _ := bot.Send(msg)
		CurrentShownLatestListID = sended.MessageID
	} else if count > 5 {
		SnippetSSSKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			row,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Показать еще", "shownextlatest"),
			),
		)
		msg.ParseMode = "HTML"
		msg.ReplyMarkup = SnippetSSSKeyboard
		sended, _ := bot.Send(msg)
		CurrentShownLatestListID = sended.MessageID
	} else {
		SnippetSSSKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			row,
		)
		msg.ParseMode = "HTML"
		msg.ReplyMarkup = SnippetSSSKeyboard
		sended, _ := bot.Send(msg)
		CurrentShownLatestListID = sended.MessageID
	}
}
func deleteSnippet(bot *tgbotapi.BotAPI, chatID int64, data *JsonWithCommandAndData) error {
	msg := tgbotapi.NewMessage(chatID, "")
	err := database.Snippets.Delete(data.ID)
	if err != nil {
		loggers.Logger.Println(err)
		return err
	}
	delRequest := tgbotapi.NewDeleteMessage(chatID, ShownSnippetMessages[data.ID])
	_, err = bot.Send(delRequest)
	if err != nil {
		loggers.Logger.Println(err)
		panic(err)
		return err
	}
	msg.Text = "Запись удалена, хорошего пивапрепровождения)"
	bot.Send(msg)
	stc := tgbotapi.NewStickerShare(chatID, "CAACAgIAAxkBAAEm_lBlMt3rvr1IsOSquJ4rqocf18MnhQACShUAAm_W4EpNHK9Mt5aCDjAE")
	bot.Send(stc)
	var count int64
	database.Snippets.DB.Table("snippets").Count(&count)
	if count == int64(ChatCurrStack[chatID]) {
		ChatCurrStack[chatID] -= 5
	}
	ReshowList(bot, chatID)
	return nil
}
func unboxSnippet(bot *tgbotapi.BotAPI, chatID int64, data *JsonWithCommandAndData) error {
	msg := tgbotapi.NewMessage(chatID, "")
	snippet, err := database.Snippets.Get(data.ID)
	if err != nil {
		loggers.Logger.Println(err)
		return err
	}
	timeForSnippet := snippet.Expires.Sub(time.Now()).Truncate(time.Hour).Hours()
	msg.Text = fmt.Sprintf("<b>%s</b>\n    %s\n        Создан <b>%s</b>\n        Пропадет через <b>%.0f</b> часов\n",
		snippet.Title, snippet.Content, snippet.Created.Format("2006-02-02 15:04"), timeForSnippet)
	msg.ParseMode = "HTML"
	dataJSONUpdate, _ := json.Marshal(JsonWithCommandAndData{
		"update", data.ID,
	})
	dataJSONDelete, _ := json.Marshal(JsonWithCommandAndData{
		"delete", data.ID,
	})
	dataJSONExtend, _ := json.Marshal(JsonWithCommandAndData{
		"extend", data.ID,
	})
	dataJSONClose, _ := json.Marshal(JsonWithCommandAndData{
		"close", data.ID,
	})
	var SnippetKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Дать леща", string(dataJSONUpdate)),
			tgbotapi.NewInlineKeyboardButtonData("Сломать колени", string(dataJSONDelete)),
			tgbotapi.NewInlineKeyboardButtonData("Напоить кумысом", string(dataJSONExtend)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Закрыть", string(dataJSONClose)),
		),
	)
	msg.ReplyMarkup = SnippetKeyboard
	sended, _ := bot.Send(msg)
	ShownSnippetMessages[data.ID] = sended.MessageID
	return nil
}
func updateSnippet(bot *tgbotapi.BotAPI, chatID int64, data *JsonWithCommandAndData) error {
	ActiveSessionIsCreatingOrEditing[chatID] = true
	UserQueryChannels[chatID] = make(chan tgbotapi.Update)
	msg := tgbotapi.NewMessage(chatID, "")
	msg.Text = "Что хотели бы изменить?"
	replyToUpdateKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Заголовок", "title"),
			tgbotapi.NewInlineKeyboardButtonData("Содержание", "content")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отмена", "cancel")))
	msg.ReplyMarkup = replyToUpdateKeyboard
	sended, _ := bot.Send(msg)
	msg.ReplyMarkup = nil
	var changedString [2]string
	for update := range UserQueryChannels[chatID] {
		if update.Message != nil {
			bot.Send(tgbotapi.NewDeleteMessage(chatID, update.Message.MessageID))
			continue
		}
		callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
		switch callback.Text {
		case "title":
			changedString[0] = callback.Text
			msg.Text = "Введите новый заголовок, ниже представлен старый\n"
			snip, _ := database.Snippets.Get(data.ID)
			msg.Text += "`" + snip.Title + "`"
			msg.ParseMode = "Markdown"
			sended2, _ := bot.Send(msg)
			for update := range UserQueryChannels[chatID] {
				changedString[1] = update.Message.Text
				err := database.Snippets.Update(data.ID, update.Message.Text, "", "")
				if err != nil {
					ActiveSessionIsCreatingOrEditing[chatID] = false
					UserQueryChannels[chatID] = nil
					loggers.Logger.Println(err)
					return err
				}
				bot.Send(tgbotapi.NewDeleteMessage(chatID, update.Message.MessageID))
				break
			}
			bot.Send(tgbotapi.NewDeleteMessage(chatID, sended2.MessageID))
			ActiveSessionIsCreatingOrEditing[chatID] = false
			UserQueryChannels[chatID] = nil

		case "content":
			changedString[0] = callback.Text
			msg.Text = "Введите новое содержание, ниже представлено старое\n"
			snip, _ := database.Snippets.Get(data.ID)
			msg.Text += "`" + snip.Content + "`"
			msg.ParseMode = "Markdown"
			sended2, _ := bot.Send(msg)
			for update := range UserQueryChannels[chatID] {
				changedString[1] = update.Message.Text
				err := database.Snippets.Update(data.ID, "", update.Message.Text, "")
				if err != nil {
					ActiveSessionIsCreatingOrEditing[chatID] = false
					UserQueryChannels[chatID] = nil
					loggers.Logger.Println(err)
					return err
				}
				bot.Send(tgbotapi.NewDeleteMessage(chatID, update.Message.MessageID))
				ActiveSessionIsCreatingOrEditing[chatID] = false
				UserQueryChannels[chatID] = nil
				break
			}
			bot.Send(tgbotapi.NewDeleteMessage(chatID, sended2.MessageID))
			ActiveSessionIsCreatingOrEditing[chatID] = false
			UserQueryChannels[chatID] = nil
		case "cancel":
			bot.Send(tgbotapi.NewDeleteMessage(chatID, sended.MessageID))
			bot.Send(tgbotapi.NewDeleteMessage(chatID, ShownSnippetMessages[data.ID]))
			ActiveSessionIsCreatingOrEditing[chatID] = false
			UserQueryChannels[chatID] = nil
			return nil
		default:
			bot.Send(tgbotapi.NewDeleteMessage(chatID, sended.MessageID))
			bot.Send(tgbotapi.NewDeleteMessage(chatID, ShownSnippetMessages[data.ID]))
			ActiveSessionIsCreatingOrEditing[chatID] = false
			UserQueryChannels[chatID] = nil
			return nil
		}
		break
	}
	bot.Send(tgbotapi.NewDeleteMessage(chatID, sended.MessageID))
	bot.Send(tgbotapi.NewDeleteMessage(chatID, ShownSnippetMessages[data.ID]))
	msg.Text = fmt.Sprintf("дело сделано, %s теперь <b>%s</b>", changedString[0], changedString[1])
	msg.ParseMode = "HTML"
	bot.Send(msg)
	stc := tgbotapi.NewStickerShare(msg.ChatID, "CAACAgIAAxkBAAEm_iRlMtTw_BwNbXhXTJOXinhlgKNy6AACNhYAAkJN2Epq39-2zr8SajAE")
	bot.Send(stc)
	err := unboxSnippet(bot, chatID, data)
	if err != nil {
		loggers.Logger.Println(err)
		return err
	}
	ActiveSessionIsCreatingOrEditing[chatID] = false
	UserQueryChannels[chatID] = nil
	return nil
}
func extendSnippet(bot *tgbotapi.BotAPI, chatID int64, data *JsonWithCommandAndData) error {
	ActiveSessionIsCreatingOrEditing[chatID] = true
	UserQueryChannels[chatID] = make(chan tgbotapi.Update)
	msg := tgbotapi.NewMessage(chatID, "")
	msg.Text = "На сколько часов хотели бы продлить?"
	sended, _ := bot.Send(msg)
	for update := range UserQueryChannels[chatID] {
		_, err := strconv.Atoi(update.Message.Text)
		if err != nil {
			msg.Text = "цифру введи дебик, ноль если переобулся"
			bot.Send(msg)
			continue
		}
		err = database.Snippets.Update(data.ID, "", "", update.Message.Text)
		if err != nil {
			ActiveSessionIsCreatingOrEditing[chatID] = false
			UserQueryChannels[chatID] = nil
			loggers.Logger.Println(err)
			return err
		}
		bot.Send(tgbotapi.NewDeleteMessage(msg.ChatID, update.Message.MessageID))
		break
	}
	bot.Send(tgbotapi.NewDeleteMessage(msg.ChatID, ShownSnippetMessages[data.ID]))
	bot.Send(tgbotapi.NewDeleteMessage(msg.ChatID, sended.MessageID))
	msg.Text = "БАХНУВ КУМЫСУ"
	bot.Send(msg)
	stc := tgbotapi.NewStickerShare(msg.ChatID, "CAACAgIAAxkBAAEm_ixlMtWbhgq7BqbuZtm6it1_uNWsHQACRhcAAjGG0EpWH82tGFrubDAE")
	bot.Send(stc)
	err := unboxSnippet(bot, chatID, data)
	if err != nil {
		ActiveSessionIsCreatingOrEditing[chatID] = false
		UserQueryChannels[chatID] = nil
		loggers.Logger.Println(err)
		return err
	}
	ActiveSessionIsCreatingOrEditing[chatID] = false
	UserQueryChannels[chatID] = nil
	return nil
}
func closeSnippet(bot *tgbotapi.BotAPI, chatID int64, data *JsonWithCommandAndData) error {
	msg := tgbotapi.NewMessage(chatID, "")
	_, err := bot.Send(tgbotapi.NewDeleteMessage(msg.ChatID, ShownSnippetMessages[data.ID]))
	if err != nil {
		loggers.Logger.Println(err)
		return err
	}
	return nil
}
