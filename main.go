package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"

	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type shortner struct {
	urls map[string]string
	mu   sync.Mutex
}

var ur = shortner{urls: make(map[string]string)}
var db *sql.DB

const (
	host     = "localhost"
	user     = "postgres"
	password = "sajjad"
	dbname   = "url_shortner"
	port     = 5432
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic("Failed to connect to database")
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic("database disconnected")
	}

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.GET("/", get)
	router.POST("/", post)
	router.GET("/:shortkey", redirect)

	router.LoadHTMLGlob("./*.html")
	err = router.Run("localhost:8080")
	if err != nil {
		fmt.Println(err)
		log.Fatal("Unable to start server: ", err)
	}
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
	ur.mu.Lock()
	stmt, err := db.Prepare("INSERT INTO url (original_url, short_url) VALUES ($1, $2)")
	if err != nil {
		ur.mu.Unlock()
		log.Println("Failed to prepare query:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save URL. Please try again."})
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(form.Originalurl, shortkey)
	if err != nil {
		ur.mu.Unlock()
		log.Printf("Failed to execute query: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save url. Please try again.",
		})
		return
	}

	ur.mu.Unlock()
	c.JSON(http.StatusOK, gin.H{"short url": shortkey, "original url": form.Originalurl})
}

func redirect(c *gin.Context) {
	var originalurl string
	shortKey := c.Param("shortkey")

	err := db.QueryRow("select original_url from url where short_url=$1", shortKey).Scan(&originalurl)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		} else {
			log.Println("Failed to retrieve URL:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve URL"})
		}
		return
	}

	c.Redirect(http.StatusMovedPermanently, originalurl)
}

func generateshortkey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	rand.Seed(time.Now().UnixMicro())
	shortkey := make([]byte, keyLength)
	for i := range shortkey {
		shortkey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortkey)
}
