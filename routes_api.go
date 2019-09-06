package main

import (
	"errors"
	"fmt"
	"net/http"
	"math/big"
	"strings"

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

	c.Status(http.StatusNoContent)
}

func RouteAPILogout(c *gin.Context) {

	// Login successful

	session := sessions.Default(c)
	session.Clear()
	session.Save()

	// Done

	c.Status(http.StatusNoContent)
}

func RouteAPIPlayer(c *gin.Context) {

	// Get server status

	status, err := rcon_command("status", "hostname: +(.*?)$")
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get player list

	players, err := get_players(status)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(players) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	// Find player

	var player *PlayerSteam
	for _, p := range players {
		if p.ID == c.Param("id") {
			player = &PlayerSteam{Player: p}
			break
		}
	}

	if player == nil {
		c.Status(http.StatusNotFound)
		return
	}

	// Split Steam ID: STEAMID_X:Y:Z

	ids := strings.Replace(player.ID, "STEAMID_", "", 1)
	idv := strings.Split(ids, ":")

	y, flag := new(big.Int).SetString(idv[1], 10)
	if !flag {
		err = errors.New("Error converting steam id part Y")
		log.Errorf("Could not set big int: %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	z, flag := new(big.Int).SetString(idv[2], 10)
	if !flag {
		err = errors.New("Error converting steam id part Z")
		log.Errorf("Could not set big int: %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// SteamID64 base

	b, flag := new(big.Int).SetString("0110000100000000", 16)
	if !flag {
		err = errors.New("Error converting hex steam base")
		log.Errorf("Could not set big int: %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert SteamID to SteamID64: z*2 + y + base

	var id64 big.Int
	id64.Mul(z, big.NewInt(2))
	id64.Add(&id64, y)
	id64.Add(&id64, b)

	// Get Steam user summary

	steam, err := get_steam(id64.String());
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if steam != nil {
		player.Steam = *steam
	}

	// Done

	c.JSON(http.StatusOK, player)
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
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, players)
}

func RouteAPIPlayersBan(c *gin.Context) {
	// Bind request body

	type Info struct {
		Message string `json:"message"`
		Minutes int    `json:"minutes"`
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

	_, err = rcon_command(fmt.Sprintf("banid %d %s", info.Minutes, id), "(.*?) was banned (.*?)$")
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

	c.Status(http.StatusNoContent)
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

	c.Status(http.StatusNoContent)
}
