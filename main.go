package main

import (
	"fmt"
	"log"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

)





func main() {
	db.opendatabase()
	var err error
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.GET("/", router.get)
	router.POST("/", post)
	router.GET("/:shortkey", redirect)

	router.LoadHTMLGlob("./template/*.html")
	err = router.Run("localhost:8080")
	if err != nil {
		fmt.Println(err)
		log.Fatal("Unable to start server: ", err)
	}
}


