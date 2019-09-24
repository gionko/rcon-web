package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RouteFEBots(c *gin.Context) {
	data := gin.H{
		"site"   : config.Site,
		"logged" : true,
		"scope"  : scope(c),
		"section": "bots",
	}

	if !authorized(c) {
		c.Redirect(http.StatusFound, config.Site.URL + "/login")
		return
	}

	c.HTML(http.StatusOK, "bots", data)
}

func RouteFEIndex(c *gin.Context) {
	data := gin.H{
		"site"   : config.Site,
		"logged" : true,
		"scope"  : scope(c),
		"section": "dashboard",
	}

	if !authorized(c) {
		c.Redirect(http.StatusFound, config.Site.URL + "/login")
		return
	}

	c.HTML(http.StatusOK, "index", data)
}

func RouteFELogin(c *gin.Context) {
	data := gin.H{
		"site"  : config.Site,
		"logged": false,
	}

	c.HTML(http.StatusOK, "login", data)
}

func RouteFEMaps(c *gin.Context) {
	data := gin.H{
		"site"   : config.Site,
		"logged" : true,
		"scope"  : scope(c),
		"section": "maps",
	}

	if !authorized(c) {
		c.Redirect(http.StatusFound, config.Site.URL + "/login")
		return
	}

	c.HTML(http.StatusOK, "maps", data)
}

func RouteFEPlayer(c *gin.Context) {
	data := gin.H{
		"site"   : config.Site,
		"logged" : true,
		"scope"  : scope(c),
		"section": "players",
		"google" : config.GoogleKey,
		"player" : c.Param("id"),
	}

	if !authorized(c) {
		c.Redirect(http.StatusFound, config.Site.URL + "/login")
		return
	}

	c.HTML(http.StatusOK, "player", data)
}

func RouteFEPlayers(c *gin.Context) {
	data := gin.H{
		"site"   : config.Site,
		"logged" : true,
		"scope"  : scope(c),
		"section": "players",
	}

	if !authorized(c) {
		c.Redirect(http.StatusFound, config.Site.URL + "/login")
		return
	}

	c.HTML(http.StatusOK, "players", data)
}
