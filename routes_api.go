package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RouteAPIUsers(c *gin.Context) {
	status, err := rcon_command("status", "hostname: +(.*?)$")
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	users, err := get_users(status)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(users) == 0 {
		c.JSON(http.StatusNoContent, nil)
		return
	}

	c.JSON(http.StatusOK, users)
}

func RouteAPIUsersBan(c *gin.Context) {
	// Bind request body

	type Info struct {
		Message string `json:"message"`
		Timeout int    `json:"timeout"`
	}
	info := Info{}

	err := c.ShouldBind(&info)
	if err != nil {
		log.Errorf("Could not bind data: %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract user id

	id := c.Param("id")

	// Ban user

	_, err = rcon_command(fmt.Sprintf("banid %d %s", info.Timeout, id), "(.*?) was banned (.*?)$")
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Kick user

	_, err = rcon_command(fmt.Sprintf("kickid %s %s", id, info.Message), "(.*?) was kicked (.*?)$")
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Done

	c.JSON(http.StatusNoContent, nil)
}

func RouteAPIUsersKick(c *gin.Context) {
	// Bind request body

	type Info struct {
		Message string `json:"message"`
	}
	info := Info{}

	err := c.ShouldBind(&info)
	if err != nil {
		log.Errorf("Could not bind data: %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract user id

	id := c.Param("id")

	// Kick user

	_, err = rcon_command(fmt.Sprintf("kickid %s %s", id, info.Message), "(.*?) was kicked (.*?)$")
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Done

	c.JSON(http.StatusNoContent, nil)
}
