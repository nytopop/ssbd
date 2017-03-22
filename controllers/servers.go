package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nytopop/ssbd/models"
)

// GET /servers/list
func ServersList(db models.Handler) gin.HandlerFunc {
	// TODO
	// get the list of servers
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "servers/list.html", gin.H{})
	}
}
