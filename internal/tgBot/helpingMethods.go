package tgBot

import (
	"SnippetsTESTBYGUIDE/internal/database"
	"SnippetsTESTBYGUIDE/internal/loggers"
	"SnippetsTESTBYGUIDE/pkg/models"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"time"
)

func ClearAllShownSnippets(bot *tgbotapi.BotAPI, msg *tgbotapi.MessageConfig) {
	for _, id := range ShownSnippetMessages {
		_, err := bot.Send(tgbotapi.NewDeleteMessage(msg.ChatID, id))
		if err != nil {
			loggers.Logger.Println(err)
		}
	}
	ShownSnippetMessages = make(map[int]int)
}
func ReshowList(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "")
	_, err := bot.Send(tgbotapi.NewDeleteMessage(chatID, CurrentShownLatestListID))
	if err != nil {
		loggers.Logger.Println(err)
		msg.Text = "Error idi nax"
		bot.Send(msg)
	}
	showLatest(bot, chatID, true)
}

func CheckIfSomeSnippetGoesExpired(bot *tgbotapi.BotAPI) {
	database.Snippets.DeleteExpired()
	expiringSnippets := make([]models.Snippet, 0, 1)
	expiring := time.Now().Add(1 * time.Hour)
	database.Snippets.DB.Model(models.Snippet{}).Where("expires BETWEEN UTC_TIMESTAMP() AND ?", expiring).Find(&expiringSnippets)
	for _, snippet := range expiringSnippets {
		creator, err := database.Users.Get(snippet.CreatedBy)
		if err != nil {
			loggers.Logger.Println(err)
			continue
		}
		msg := tgbotapi.NewMessage(creator.ChatID, fmt.Sprintf("SNIPPET %s IS EXPIRING", snippet.Title))
		bot.Send(msg)
	}
}
