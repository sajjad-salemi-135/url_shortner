package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sajjad-salemi-135/url_shortner/handler"
	"github.com/sajjad-salemi-135/url_shortner/db"
)


func main() {
	db.opendatabase()
	var err error
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.GET("/", modul.get)
	router.POST("/", modul.post)
	router.GET("/:shortkey", modul.redirect)

	router.LoadHTMLGlob("./template/*.html")
	err = router.Run("localhost:8080")
	if err != nil {
		fmt.Println(err)
		log.Fatal("Unable to start server: ", err)
	}
}


