package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nytopop/ssbd/models"
)

// GET /dash/overview
func DashOverview(db models.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "dash/overview.html", gin.H{})
	}
}
