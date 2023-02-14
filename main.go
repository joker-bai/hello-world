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
	mysqlUser = "root"
	mysqlPassword = "Resico@2020#dev"
	mysqlAddress = "192.168.100.23:3306"
	mysqlDBName = "hello-world"
	fmt.Print(mysqlUser, mysqlPassword, mysqlAddress, mysqlDBName)
	fmt.Print("============")
	// init database
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlUser, mysqlPassword, mysqlAddress, mysqlDBName)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	// auto migrate
	db.AutoMigrate(&HelloWorld{})
	db.Create(&HelloWorld{Text: "Hello World"})
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
	g.POST("/", func(context *gin.Context) {
		// bind
		var hello HelloWorld
		if err := context.ShouldBind(&hello); err != nil {
			context.JSON(304, gin.H{"status": err.Error()})
		}
		fmt.Println(hello)
		// set content to mysql
		db.Model(&HelloWorld{}).Where("id=?", 1).Update("text", hello.Text)
		context.Redirect(http.StatusFound, "/")
	})

	if err := g.Run(); err != nil {
		log.Fatalln("start failed!")
	}
}
