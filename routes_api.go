package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/madcitygg/rcon"
)

func RouteAPILogin(c *gin.Context) {
	// Bind request body

	type Login struct {
		Password string `json:"password"`
	}
	login := Login{}

	err := c.ShouldBind(&login)
	if err != nil {
		log.Errorf("Could not bind data: %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Try to authenticate with provided password

	req, err := rcon.Dial(fmt.Sprintf("%s:%d", config.ServerAddress, config.ServerPort))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer req.Close()

	err = req.Authenticate(login.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Login successful

	session := sessions.Default(c)
	session.Set("logged", true)
	session.Set("password", login.Password)
	session.Save()

	// Done

	c.JSON(http.StatusNoContent, nil)
}

func RouteAPILogout(c *gin.Context) {

	// Login successful

	session := sessions.Default(c)
	session.Clear()
	session.Save()

	// Done

	c.JSON(http.StatusNoContent, nil)
}

func RouteAPIPlayers(c *gin.Context) {
	status, err := rcon_command("status", "hostname: +(.*?)$")
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	players, err := get_players(status)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(players) == 0 {
		c.JSON(http.StatusNoContent, nil)
		return
	}

	c.JSON(http.StatusOK, players)
}

func RouteAPIPlayersBan(c *gin.Context) {
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

func RouteAPIPlayersKick(c *gin.Context) {
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
