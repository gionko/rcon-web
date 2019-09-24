package main

import (
	"errors"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func authorized(c *gin.Context) bool {
	session := sessions.Default(c)
	v := session.Get("logged")
	if v != nil {
		return v.(bool)
	}

	return false
}

func scope(c *gin.Context) string {
	session := sessions.Default(c)
	v := session.Get("scope")
	if v != nil {
		return v.(string)
	}

	return ""
}

func username(c *gin.Context) (string, error) {

	// Extract user name from session

	session := sessions.Default(c)
	v := session.Get("name")
	if v == nil {
		err := errors.New("Could not extract name from session data")
		log.Errorf("Unauthorized API action: %+v", err)
		return "", err
	}
	name := v.(string)

	// Done

	return name, nil
}
