package database

import (
	"SnippetsTESTBYGUIDE/internal/loggers"
	"SnippetsTESTBYGUIDE/pkg/models"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type SnippetModel struct {
	DB *gorm.DB
}

var Snippets SnippetModel

func InitSnippetModel() error {
	err := DbGORM.AutoMigrate(&models.Snippet{})
	Snippets.DB = DbGORM
	return err
}
func (m *SnippetModel) Insert(userID int, title, content, expires string) (int, error) {
	exp, _ := strconv.Atoi(expires)
	m.DB.Model(&models.Snippet{}).Create(&models.Snippet{
		Title:     title,
		Content:   content,
		Created:   time.Now(),
		Expires:   time.Now().Add((time.Duration(exp)) * time.Hour),
		CreatedBy: userID,
	})
	return 0, nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	var found *models.Snippet
	if err := m.DB.Model(&models.Snippet{}).First(&found, id).Error; err != nil {
		loggers.Logger.Println(err)
		return nil, models.ErrNoRecordedSnippet
	}
	return found, nil
}
func (m *SnippetModel) Update(id int, title, content, expires string) error {
	var exp int
	if expires != "0" {
		exp, _ = strconv.Atoi(expires)
	} else {
		exp = 0
	}
	modelToUpdate := models.Snippet{}
	m.DB.Model(&models.Snippet{}).Find(&models.Snippet{}, id).First(&modelToUpdate)
	if title != "" {
		modelToUpdate.Title = title
	}
	if content != "" {
		modelToUpdate.Content = content
	}
	if exp != 0 {
		modelToUpdate.Expires = modelToUpdate.Expires.Add(time.Duration(exp) * time.Hour)
	}
	err := m.DB.Save(&modelToUpdate).Error
	if err != nil {
		loggers.Logger.Println(err)
		return err
	}
	return nil
}

func (m *SnippetModel) Latest(offset, limit, createdBy int) ([]*models.Snippet, error) {
	var found []*models.Snippet
	err := m.DB.Model(&models.Snippet{}).Order("id desc").Offset(offset).Limit(limit).Where("expires > UTC_TIMESTAMP() AND created_by = ?", createdBy).Find(&found).Error
	if err != nil {
		loggers.Logger.Println(err)
		return nil, err
	}
	return found, nil
}

func (m *SnippetModel) Delete(id int) error {
	err := m.DB.Model(&models.Snippet{}).Delete(&models.Snippet{}, id).Error
	if err != nil {
		loggers.Logger.Println(err)
		return err
	}
	return nil
}
func (m *SnippetModel) DeleteExpired() {
	m.DB.Model(&models.Snippet{}).Where("expires < UTC_TIMESTAMP()").Delete(&models.Snippet{})
}
