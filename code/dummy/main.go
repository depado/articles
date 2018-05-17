package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"alive": true})
	})

	if err := r.Run("127.0.0.1:8080"); err != nil {
		logrus.WithError(err).Fatal("Couldn't listen")
	}
}
