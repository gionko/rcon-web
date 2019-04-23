package main

import (
	"net/http"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func multi_render() multitemplate.Render {
	r := multitemplate.New()
	r.AddFromFiles("index", "templates/base.tmpl", "templates/index.tmpl")
	r.AddFromFiles("login", "templates/base.tmpl", "templates/login.tmpl")
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
