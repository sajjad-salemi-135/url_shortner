package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sajjad-salemi-135/url_shortner/db"
	"github.com/sajjad-salemi-135/url_shortner/handler"
)


func main() {
	db.Opendatabase()
	var err error
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.LoadHTMLGlob("./template/*.html")

	router.GET("/", handler.Get)
	router.POST("/", handler.Post)
	router.GET("/:shortkey", handler.Redirect)


	err = router.Run("localhost:8080")
	if err != nil {
		fmt.Println(err)
		log.Fatal("Unable to start server: ", err)
	}
}


