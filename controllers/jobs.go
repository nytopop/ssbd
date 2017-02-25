package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nytopop/ssbd/models"
)

// GET /jobs/list
func JobsList(db models.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "jobs/list.html", gin.H{})
	}
}
