package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/Depado/articles/code/dialogflow/cocktail"
	"github.com/gin-gonic/gin"
	df "github.com/leboncoin/dialogflow-go-webhook"
)

func search(c *gin.Context, dfr *df.Request) {
	c.JSON(http.StatusOK, gin.H{})
}

func random(c *gin.Context, dfr *df.Request) {
	var err error
	var d *cocktail.FullDrink

	if d, err = cocktail.C.GetRandomDrink(); err != nil {
		logrus.WithError(err).Error("Coudln't get random drink")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	out := fmt.Sprintf("I found that cocktail : %s", d.StrDrink)
	dff := &df.Fulfillment{
		FulfillmentMessages: df.Messages{
			{RichMessage: df.Text{Text: []string{out}}},
			df.ForGoogle(df.SingleSimpleResponse(out, out)),
		},
	}
	c.JSON(http.StatusOK, dff)
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
