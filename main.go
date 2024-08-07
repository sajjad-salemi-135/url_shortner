package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sajjad-salemi-135/url_shortner/handler"
	"github.com/sajjad-salemi-135/url_shortner/db"
)


func main() {
	db.opendatabase()
	var err error
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.GET("/", handler.get)
	router.POST("/", handler.post)
	router.GET("/:shortkey", handler.redirect)

	router.LoadHTMLGlob("./template/*.html")
	err = router.Run("localhost:8080")
	if err != nil {
		fmt.Println(err)
		log.Fatal("Unable to start server: ", err)
	}
}


