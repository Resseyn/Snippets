package database

import (
	"SnippetsTESTBYGUIDE/internal/loggers"
	"SnippetsTESTBYGUIDE/pkg/models"
	"gorm.io/gorm"
)

type UserModel struct {
	DB *gorm.DB
}

var Users UserModel

func InitUserModel() error {
	Users.DB = DbGORM
	return nil
}

func (m *UserModel) Insert(userID int, chatID int64) (int, error) {
	err := m.DB.Model(&models.User{}).Create(&models.User{
		UserID: userID,
		ChatID: chatID,
	}).Error
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func (m *UserModel) Get(userID int) (*models.User, error) {
	var found *models.User
	if err := m.DB.Model(&models.User{}).First(&found, userID).Error; err != nil {
		loggers.Logger.Println(err)
		return nil, models.ErrNoRecordedSnippet
	}
	return found, nil
}
func (m *UserModel) Update(userID int, chatID int64) error {

	modelToUpdate := models.User{}
	m.DB.Model(&models.User{}).Find(&models.User{}, userID).First(&modelToUpdate)
	modelToUpdate.ChatID = chatID
	err := m.DB.Save(&modelToUpdate).Error
	if err != nil {
		loggers.Logger.Println(err)
		return err
	}
	return nil
}
