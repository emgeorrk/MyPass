package main

import (
	"encoding/base64"
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
	db = new(dataBase)
)

func auth(c *gin.Context) {
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

func abortWithStatusAndError(status int, err error, c *gin.Context) {
	log.Println(err)
	c.String(status, err.Error())
	c.Abort()
}

func actionHandler(c *gin.Context) {
	action := c.GetHeader("Action")

	switch action {
	case "getBase":
		getBasePrep(c)
	case "addElem":
		addElemPrep(c)
	case "editElem":
		editElemPrep(c)
	case "removeElem":
		removeElemPrep(c)
	default:
		c.String(http.StatusBadRequest, "Action not supported: %s\n", action)
	}
}

func getBasePrep(c *gin.Context) {
	base, status, err := db.getBase()
	if err != nil {
		abortWithStatusAndError(status, err, c)
		return
	}

	c.JSONP(http.StatusOK, base)
	c.String(http.StatusOK, "\n")
}

func addElemPrep(c *gin.Context) {
	newElem := Element{}
	if status, err := newElem.getElem("Element", c); err != nil {
		abortWithStatusAndError(status, err, c)
		return
	}

	if status, err := db.addElem(newElem); err != nil {
		abortWithStatusAndError(status, err, c)
		return
	}

	c.String(http.StatusOK, "Element added successfully\n")
	getBasePrep(c)
}

func editElemPrep(c *gin.Context) {
	oldElem := Element{}
	if status, err := oldElem.getElem("oldElement", c); err != nil {
		abortWithStatusAndError(status, err, c)
		return
	}

	newElem := Element{}
	if status, err := newElem.getElem("newElement", c); err != nil {
		abortWithStatusAndError(status, err, c)
		return
	}

	status, err := db.editElem(oldElem, newElem)
	if err != nil {
		abortWithStatusAndError(status, err, c)
		return
	}

	c.String(http.StatusOK, "Element edited successfully\n")
	getBasePrep(c)
}

func removeElemPrep(c *gin.Context) {
	elem := Element{}
	if status, err := elem.getElem("Element", c); err != nil {
		abortWithStatusAndError(status, err, c)
		return
	}

	if status, err := db.removeElem(elem); err != nil {
		abortWithStatusAndError(status, err, c)
		return
	}

	c.String(http.StatusOK, "Element removed successfully\n")
	getBasePrep(c)
}

func main() {
	var err error
	err = db.connectDB()
	if err != nil {
		log.Fatalln("Error connecting PostgreSQL: ", err)
	}
	defer db.postgres.Close()

	r := gin.Default()

	r.Handle("GET", "/ping", auth, actionHandler)

	err = r.Run()
	if err != nil {
		log.Fatalln("Error launching server: ", err)
	}
}
