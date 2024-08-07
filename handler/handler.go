package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sajjad-salemi-135/url_shortner/db"
)

func Get(c *gin.Context) {
	c.HTML(http.StatusOK, "get.html", gin.H{"title": "main web"})
}

func Post(c *gin.Context) {
	var form struct {
		Originalurl string `form:"original_url"`
	}
	if err := c.ShouldBind(&form); err != nil {
		c.String(404, "bad req", gin.H{"error": "All fields are required and must be valid."})
	}

	shortkey := Generateshortkey()

	db.Postdb(shortkey, form.Originalurl, c)

	c.JSON(http.StatusOK, gin.H{"short url": shortkey, "original url": form.Originalurl})
}

func Redirect(c *gin.Context) {
	var originalurl string
	shortKey := c.Param("shortkey")

	originalurl = db.Redirectdb(shortKey, c)

	c.Redirect(http.StatusMovedPermanently, originalurl)
}
