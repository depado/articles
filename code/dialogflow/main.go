package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	df "github.com/leboncoin/dialogflow-go-webhook"
)

func search(c *gin.Context, dfr *df.Request) {
	c.JSON(http.StatusOK, gin.H{})
}

func random(c *gin.Context, dfr *df.Request) {
	c.JSON(http.StatusOK, gin.H{})
}

func webhook(c *gin.Context) {
	var err error
	var dfr *df.Request

	if err = c.BindJSON(&dfr); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	switch dfr.QueryResult.Action {
	case "search":
		log.Println("Search action detected")
		search(c, dfr)
	case "random":
		log.Println("Random action detected")
		random(c, dfr)
	default:
		log.Println("Unknown action")
		c.AbortWithStatus(http.StatusNotFound)
	}
}

func main() {
	r := gin.Default()
	r.POST("/webhook", webhook)
	if err := r.Run("127.0.0.1:8001"); err != nil {
		panic(err)
	}
}
