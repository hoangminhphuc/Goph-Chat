package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/hoangminhphuc/goph-chat/migration"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	
	dbUri := os.Getenv("MYSQL_GORM_DB_URI")
	db, err := gorm.Open(mysql.Open(dbUri), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.Debug()

	err = migration.AutoMigrate(db)

	router := gin.Default()



	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})

	router.Run(":8080")
}