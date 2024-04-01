package main

import (
	"encoding/base64"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

var (
	username = "admin"
	password = "admin"
) // admin:admin YWRtaW46YWRtaW4=

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

func abortWithStatus(statusCode int, err error, c *gin.Context) {
	log.Println(err)
	c.JSON(statusCode, gin.H{"error": err.Error()})
	c.Abort()
}

func getBasePrep(c *gin.Context) {
	base, statusCode, err := db.getBase()
	if err != nil {
		abortWithStatus(statusCode, err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"passwords": base.Rows})
}

func addElemPrep(c *gin.Context) {
	newElem := element{}
	if statusCode, err := newElem.getElem("Element", c); err != nil {
		abortWithStatus(statusCode, err, c)
		return
	}

	if statusCode, err := db.addElem(newElem); err != nil {
		abortWithStatus(statusCode, err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func editElemPrep(c *gin.Context) {
	oldElem := element{}
	if statusCode, err := oldElem.getElem("oldElement", c); err != nil {
		abortWithStatus(statusCode, err, c)
		return
	}

	newElem := element{}
	if statusCode, err := newElem.getElem("newElement", c); err != nil {
		abortWithStatus(statusCode, err, c)
		return
	}

	statusCode, err := db.editElem(oldElem, newElem)
	if err != nil {
		abortWithStatus(statusCode, err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func removeElemPrep(c *gin.Context) {
	elem := element{}
	if statusCode, err := elem.getElem("Element", c); err != nil {
		abortWithStatus(statusCode, err, c)
		return
	}

	if statusCode, err := db.removeElem(elem); err != nil {
		abortWithStatus(statusCode, err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func editAuth(c *gin.Context) {
	// TO BE CONTINUED
}

func main() {
	var err error
	err = db.connectDB()
	if err != nil {
		log.Fatalln("Error connecting PostgreSQL: ", err)
	}
	defer db.postgres.Close()

	r := gin.Default()

	//r.ForwardedByClientIP = true
	//r.SetTrustedProxies([]string{"1.1.1.1"})

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"*"}

	r.Use(cors.New(config))

	r.OPTIONS("/", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusOK)
	})
	r.GET("/", auth, getBasePrep)
	r.PUT("/", auth, editElemPrep)
	r.POST("/", auth, addElemPrep)
	r.DELETE("/", auth, removeElemPrep)
	r.PATCH("/", auth, editAuth)

	if err := r.Run(":56821"); err != nil {
		log.Fatalln("Error launching server:", err)
	}
}
