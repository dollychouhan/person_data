package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func registerRoutes(router *gin.Engine, db *sql.DB) {
	router.GET("/person/:person_id/info", func(c *gin.Context) {
		personId := c.Param("person_id")
		person, err := GetPersonInfo(db, personId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, person)
	})
}
