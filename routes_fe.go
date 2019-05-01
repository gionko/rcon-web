package main

import (
	"net/http"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func multi_render() multitemplate.Render {
	r := multitemplate.New()
	r.AddFromFiles("bots",    "templates/base.tmpl", "templates/bots.tmpl")
	r.AddFromFiles("index",   "templates/base.tmpl", "templates/index.tmpl")
	r.AddFromFiles("login",   "templates/base.tmpl", "templates/login.tmpl")
	r.AddFromFiles("maps",    "templates/base.tmpl", "templates/maps.tmpl")
	r.AddFromFiles("players", "templates/base.tmpl", "templates/players.tmpl")
	return r
}

func authorized(c *gin.Context) bool {
	session := sessions.Default(c)
	v := session.Get("logged")
	if v != nil {
		return v.(bool)
	}

	return false
}

func RouteFEBots(c *gin.Context) {
	data := gin.H{
		"site"   : config.Site,
		"logged" : true,
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
		"section": "maps",
	}

	if !authorized(c) {
		c.Redirect(http.StatusFound, config.Site.URL + "/login")
		return
	}

	c.HTML(http.StatusOK, "maps", data)
}

func RouteFEPlayers(c *gin.Context) {
	data := gin.H{
		"site"   : config.Site,
		"logged" : true,
		"section": "players",
	}

	if !authorized(c) {
		c.Redirect(http.StatusFound, config.Site.URL + "/login")
		return
	}

	c.HTML(http.StatusOK, "players", data)
}
