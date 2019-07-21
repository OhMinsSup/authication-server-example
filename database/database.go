package database

import (
	"fmt"
	"github.com/OhMinsSup/lafu-server/database/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
)

// Initialize 데이터베이스 초기화
func Initialize() (*gorm.DB, error) {
	dbConfig := os.Getenv("DB_CONFIG")
	db, err := gorm.Open("postgres", dbConfig)

	// Logs SQL
	db.LogMode(true)
	db.Set("gorm:table_options", "charset=utf8")
	// created uuid
	db.Callback().Create().Before("gorm:create").Register("my_plugin:before_create", models.BeforeCreateUUID)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to database")
	models.Migrate(db)

	return db, err
}
