package db

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	
)

var db *sql.DB


const (
	host     = "localhost"
	user     = "postgres"
	password = "sajjad"
	dbname   = "url_shortner"
	port     = 5432
)

func opendatabase() {
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
}

func redirectdb(shortKey string,originalurl string,c *gin.Context){
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
}

func postdb(shortkey string, originalurl string, c *gin.Context){
	ur.mu.Lock()
	stmt, err := db.Prepare("INSERT INTO url (original_url, short_url) VALUES ($1, $2)")
	if err != nil {
		ur.mu.Unlock()
		log.Println("Failed to prepare query:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save URL. Please try again."})
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(originalurl, shortkey)
	if err != nil {
		ur.mu.Unlock()
		log.Printf("Failed to execute query: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save url. Please try again.",
		})
		return
	}

	ur.mu.Unlock()
}