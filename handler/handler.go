package router

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sajjad-salemi-135/url_shortner/db"
)

var ur = shortner{urls: make(map[string]string)}

type shortner struct {
	urls map[string]string
	mu   sync.Mutex
}

func get(c *gin.Context) {
	c.HTML(http.StatusOK, "get.html", gin.H{"title": "main web"})
}

func post(c *gin.Context) {
	var form struct {
		Originalurl string `form:"original_url"`
	}
	if err := c.ShouldBind(&form); err != nil {
		c.String(404, "bad req", gin.H{"error": "All fields are required and must be valid."})
	}

	shortkey := generateshortkey()

	db.postdb(shortkey,form.Originalurl,c)
	
	c.JSON(http.StatusOK, gin.H{"short url": shortkey, "original url": form.Originalurl})
}

func redirect(c *gin.Context) {
	var originalurl string
	shortKey := c.Param("shortkey")

	db.redirectdb(shortKey,originalurl,c)

	c.Redirect(http.StatusMovedPermanently, originalurl)
}

