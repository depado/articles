package main

import (
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"

	"github.com/Depado/articles/code/dialogflow/cocktail"
	"github.com/gin-gonic/gin"
	df "github.com/leboncoin/dialogflow-go-webhook"
)

func cardFromDrink(d *cocktail.FullDrink) df.BasicCard {
	card := df.BasicCard{
		Title:         d.StrDrink,
		FormattedText: d.StrInstructions,
		Image: &df.Image{
			ImageURI: d.StrDrinkThumb,
		},
	}
	return card
}

type searchParams struct {
	Alcohol   string `json:"alcohol"`
	DrinkType string `json:"drink-type"`
	Name      string `json:"name"`
}

func search(c *gin.Context, dfr *df.Request) {
	var err error
	var p searchParams

	if err = dfr.GetParams(&p); err != nil {
		logrus.WithError(err).Error("Couldn't get parameters")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	spew.Dump(p)

	c.JSON(http.StatusOK, gin.H{})
}

func specify(c *gin.Context, dfr *df.Request) {
	var err error
	var p searchParams

	if err = dfr.GetContext("Search-followup", &p); err != nil {
		logrus.WithError(err).Error("Couldn't get parameters")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	spew.Dump(p)

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
			df.ForGoogle(cardFromDrink(d)),
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
	clog := logrus.WithField("action", dfr.QueryResult.Action)

	switch dfr.QueryResult.Action {
	case "search":
		clog.Info("Detected")
		search(c, dfr)
	case "random":
		clog.Info("Detected")
		random(c, dfr)
	case "search.specify":
		clog.Info("Detected")
	default:
		clog.Warn("Unknown")
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
