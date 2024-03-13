package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

const (
	username = "egormerk"
	password = "admin"
)

var (
	db *sql.DB
)

func Auth(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		return
	}

	if len(authHeader) < 6 || authHeader[:6] != "Basic " {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Malformed Authorization header"})
		return
	}

	basicAuthData := authHeader[6:]
	decodedData, err := base64.StdEncoding.DecodeString(basicAuthData)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Failed to decode Authorization header"})
		return
	}

	credentials := strings.Split(string(decodedData), ":")
	if len(credentials) != 2 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
		return
	}

	if credentials[0] != username || credentials[1] != password {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	c.Next()
}

func foo2(c *gin.Context) {
	//c.String(http.StatusOK, "Authorization successful\n")
	c.Next()
}

func foo3(c *gin.Context) {
	action := c.GetHeader("Action")

	switch action {
	case "getBase":
		base, status, err := getBase(db)
		if err != nil {
			log.Println(err)
			c.String(status, err.Error())
			c.Abort()
			return
		}

		c.JSONP(http.StatusOK, base)
		c.String(http.StatusOK, "\n")

	case "addElem":
		jsonNewElem := c.GetHeader("Element")
		newElem := Element{}
		err := json.Unmarshal([]byte(jsonNewElem), &newElem)
		if err != nil {
			log.Println(err)
			c.String(http.StatusInternalServerError, err.Error())
			c.Abort()
			return
		}

		status, err := addElem(db, newElem.Title, newElem.Login, newElem.Password)
		if err != nil {
			log.Println(err)
			c.String(status, err.Error())
			c.Abort()
			return
		}

		c.String(http.StatusOK, "Element added successfully\n")
	default:
		c.String(http.StatusBadRequest, "Action not supported\n")
	}
}

func main() {
	var err error
	db, err = connectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := gin.Default()
	r.Handle("GET", "/ping", Auth, foo2, foo3)
	err = r.Run()
	if err != nil {
		log.Fatal("Error launching server")
	}
}
