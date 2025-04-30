package gormdiaclects

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func MySQLConnection(uri string) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(uri), &gorm.Config{})
}