package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nytopop/ssbd/models"
)

type AddServerForm struct {
	Name    string
	Address string
	Port    int
}

// GET /servers/list
func ServersList(db models.Handler) gin.HandlerFunc {
	// TODO
	// get the list of servers
	return func(c *gin.Context) {
		srvs, err := db.GetServers()
		if err != nil {
			RenderErr(c, err)
			return
		}

		c.HTML(http.StatusOK, "servers/list.html", gin.H{
			"Servers": srvs,
		})
	}
}
