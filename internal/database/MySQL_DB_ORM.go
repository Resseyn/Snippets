package database

import (
	"database/sql"
	"flag"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DbGORM *gorm.DB

func InitMySqlORMdatabase() error {
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "Название MySQL источника данных")
	flag.Parse()

	db, err := openDB(*dsn)
	if err != nil {
		return err
	}
	DbGORM, err = gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
