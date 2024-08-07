package db

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB
var ur = shortner{urls: make(map[string]string)}

type shortner struct {
	urls map[string]string
	mu   sync.Mutex
}

const (
	host     = "localhost"
	user     = "postgres"
	password = "sajjad"
	dbname   = "url_shortner"
	port     = 5432
)

func Opendatabase() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic("Failed to connect to database")
	}

	err = db.Ping()
	if err != nil {
		panic("database disconnected")
	}
}

func Closedb() {
	db.Close()
}

func Redirectdb(shortKey string, c *gin.Context) string {
	var originalurl string
	err := db.QueryRow("select original_url from url where short_url=$1", shortKey).Scan(&originalurl)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		} else {
			log.Println("Failed to retrieve URL:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve URL"})
		}
		return ""
	}
	return originalurl
}

func Postdb(shortkey string, originalurl string, c *gin.Context) {
	ur.mu.Lock()
	defer ur.mu.Unlock()
	stmt, err := db.Prepare("INSERT INTO url (original_url, short_url) VALUES ($1, $2)")
	if err != nil {
		log.Println("Failed to prepare query:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "insert Failed to save URL. Please try again."})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(originalurl, shortkey)
	if err != nil {
		log.Printf("Failed to execute query: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "excute Failed to save url. Please try again.",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "URL saved successfully."})
}
