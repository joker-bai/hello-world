package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

var (
	db  *gorm.DB
	err error
)

// HelloWorld Model Struct
type HelloWorld struct {
	ID   int8   `gorm:"primaryKey,autoIncrement"`
	Text string `json:"text" form:"text"`
}

// init
func init() {
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	mysqlAddress := os.Getenv("MYSQL_ADDRESS")
	mysqlDBName := os.Getenv("MYSQL_DBNAME")
	content := os.Getenv("CONTENT")
	fmt.Print(mysqlUser, mysqlPassword, mysqlAddress, mysqlDBName)
	// init database
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlUser, mysqlPassword, mysqlAddress, mysqlDBName)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	// auto migrate
	db.AutoMigrate(&HelloWorld{})
	db.Create(&HelloWorld{Text: content})
}

// main
func main() {
	g := gin.New()
	g.GET("/", func(context *gin.Context) {
		// get content from mysql
		var hello HelloWorld
		db.First(&hello, "1")
		g.LoadHTMLFiles("index.html")
		context.HTML(http.StatusOK, "index.html", hello)
	})

	if err := g.Run(); err != nil {
		log.Fatalln("start failed!")
	}
}
