package main

import (
	"errors"
	"fmt"
	"net/http"
	"math/big"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func RouteAPIBots(c *gin.Context) {

	// Check if authorized

	if !authorized(c) {
		c.Status(http.StatusUnauthorized)
		return
	}

	// Change number of bots

	_, err := rcon_command("ins_bot_count_checkpoint " + c.Param("id"), ".*?ins_bot_count_checkpoint *(\\d*).*$")
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Log action

	session := sessions.Default(c)
	v := session.Get("name")
	if v == nil {
		err = errors.New("Could not extract name from session data")
		log.Errorf("Unauthorized API action: %+v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}
	name := v.(string)

	log.Infof("User %s changed amount of bots to %s", name, c.Param("id"))

	// Done

	c.Status(http.StatusNoContent)
}

func RouteAPIDamage(c *gin.Context) {

	// Check if authorized

	if !authorized(c) {
		c.Status(http.StatusUnauthorized)
		return
	}

	// Change bot damage

	_, err := rcon_command("bot_damage " + c.Param("id"), ".*?bot_damage *(\\d*).*$")
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Log action

	session := sessions.Default(c)
	v := session.Get("name")
	if v == nil {
		err = errors.New("Could not extract name from session data")
		log.Errorf("Unauthorized API action: %+v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}
	name := v.(string)

	log.Infof("User %s changed bot damage to %s", name, c.Param("id"))

	// Done

	c.Status(http.StatusNoContent)
}

func RouteAPIDifficulty(c *gin.Context) {

	// Check if authorized

	if !authorized(c) {
		c.Status(http.StatusUnauthorized)
		return
	}

	// Change bot damage

	_, err := rcon_command("ins_bot_difficulty " + c.Param("id"), ".*?ins_bot_difficulty *(\\d*).*$")
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Log action

	session := sessions.Default(c)
	v := session.Get("name")
	if v == nil {
		err = errors.New("Could not extract name from session data")
		log.Errorf("Unauthorized API action: %+v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}
	name := v.(string)

	log.Infof("User %s changed bot difficulty to %s", name, c.Param("id"))

	// Done

	c.Status(http.StatusNoContent)
}

func RouteAPILogin(c *gin.Context) {

	// Bind request body

	type Login struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	login := Login{}

	err := c.ShouldBind(&login)
	if err != nil {
		log.Errorf("Could not bind data: %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user in configuration

	found := false
	var password string
	var scope string

	for _, u := range config.Users {
		if u.Name == login.Name {
			found = true
			password = u.Password
			scope = u.Scope
			break
		}
	}

	if !found {
		c.Status(http.StatusNotFound)
		return
	}

	// Match passwords

	if login.Password != password {
		c.Status(http.StatusUnauthorized)
		return
	}

	// Login successful

	session := sessions.Default(c)
	session.Set("logged", true)
	session.Set("name",   login.Name)
	session.Set("scope",  scope)
	session.Save()

	// Log action

	log.Infof("User %s (%s) logged in", login.Name, scope)

	// Done

	c.Status(http.StatusNoContent)
}

func RouteAPILogout(c *gin.Context) {

	// Log action

	var name string

	session := sessions.Default(c)
	v := session.Get("name")
	if v != nil {
		name = v.(string)
	} else {
		err := errors.New("Could not extract name from session data")
		log.Errorf("Unauthorized API action: %+v", err)
		name = "<unknown>"
	}

	log.Infof("User %s logged out", name)

	// Login successful

	session.Clear()
	session.Save()

	// Done

	c.Status(http.StatusNoContent)
}

func RouteAPIMap(c *gin.Context) {

	// Check if authorized

	if !authorized(c) {
		c.Status(http.StatusUnauthorized)
		return
	}

	// Extract map id

	id := c.Param("id")

	// Log action

	session := sessions.Default(c)
	v := session.Get("name")
	if v == nil {
		err := errors.New("Could not extract name from session data")
		log.Errorf("Unauthorized API action: %+v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}
	name := v.(string)

	log.Infof("User %s changed map to %s", name, id)

	// Change map

	/* TODO: make map type configurable, currently hardcoded to 'checkpoint' */

	_, err := rcon_command("changelevel " + id + " checkpoint", "")
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Done

	c.Status(http.StatusNoContent)
}

func RouteAPIMaps(c *gin.Context) {

	// Check if authorized

	if !authorized(c) {
		c.Status(http.StatusUnauthorized)
		return
	}

	// Get map list from server

	/* TODO: make mask configurable */

	reply, err := rcon_command("maps coop", "-------------")
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get map list

	maps, err := get_maps(reply)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(maps) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	// Done

	c.JSON(http.StatusOK, maps)
}

func RouteAPIPlayer(c *gin.Context) {

	// Check if authorized

	if !authorized(c) {
		c.Status(http.StatusUnauthorized)
		return
	}

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

	// Check if authorized

	if !authorized(c) {
		c.Status(http.StatusUnauthorized)
		return
	}

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

	// Done

	c.JSON(http.StatusOK, players)
}

func RouteAPIPlayersBan(c *gin.Context) {

	// Check if authorized

	if !authorized(c) {
		c.Status(http.StatusUnauthorized)
		return
	}

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

	// Check permissions

	if scope(c) != "admin" && info.Minutes > 1440 {
		err = errors.New("Only admins has permission to ban users for more than 1 day")
		log.Errorf("Unauthorized API action: %+v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Extract user id

	id := c.Param("id")

	// Log action

	session := sessions.Default(c)
	v := session.Get("name")
	if v == nil {
		err = errors.New("Could not extract name from session data")
		log.Errorf("Unauthorized API action: %+v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}
	name := v.(string)

	log.Infof("User %s banned player %s for %d minutes, %s", name, id, info.Minutes, info.Message)

	// Ban user

	_, err = rcon_command(fmt.Sprintf("banid %d %s", info.Minutes, id), "(.*?)was banned(.*?)$")
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Kick user

	_, err = rcon_command(fmt.Sprintf("kickid %s %s", id, info.Message), "(.*?)(Kicked by Console)|(not found)(.*?)$")
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Done

	c.Status(http.StatusNoContent)
}

func RouteAPIPlayersKick(c *gin.Context) {

	// Check if authorized

	if !authorized(c) {
		c.Status(http.StatusUnauthorized)
		return
	}

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

	// Log action

	session := sessions.Default(c)
	v := session.Get("name")
	if v == nil {
		err = errors.New("Could not extract name from session data")
		log.Errorf("Unauthorized API action: %+v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}
	name := v.(string)

	log.Infof("User %s kicked player %s, %s", name, id, info.Message)

	// Kick user

	_, err = rcon_command(fmt.Sprintf("kickid %s %s", id, info.Message), "(.*?)(Kicked by Console)|(not found)(.*?)$")
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Done

	c.Status(http.StatusNoContent)
}

func RouteAPIStatus(c *gin.Context) {

	// Check if authorized

	if !authorized(c) {
		c.Status(http.StatusUnauthorized)
		return
	}

	// Get server status

	reply, err := rcon_command("status", "hostname: +(.*?)$")
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract status data

	status, err := get_status(reply)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Done

	c.JSON(http.StatusOK, status)
}
