package models

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/pborman/uuid"
)

// Migrate automigrates schema using ORM
func Migrate(db *gorm.DB) {
	db.AutoMigrate(&User{}, &AuthToken{})
	db.Model(&AuthToken{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	fmt.Println("Auto Migration has beed processed")
}

// BeforeCreateUUID RDBMS의 UUID를 자동으로 생성해주는 함수
func BeforeCreateUUID(scope *gorm.Scope) {
	reflectValue := reflect.Indirect(reflect.ValueOf(scope.Value))
	if strings.Contains(string(reflectValue.Type().Field(0).Tag), "uuid") {
		uuid.SetClockSequence(-1)
		scope.SetColumn("id", uuid.NewUUID().String())
	}
}
