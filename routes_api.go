package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RouteAPIUsers(c *gin.Context) {
	status, err := get_status()
	if err != nil {
		log.Errorf("Could not obtain `RCON status response: %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	users, err := get_users(status)
	if err != nil {
		log.Errorf("Could not obtain player scores: %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(users) == 0 {
		c.JSON(http.StatusNoContent, nil)
		return
	}

	c.JSON(http.StatusOK, users)
}
