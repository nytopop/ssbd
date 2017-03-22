package controllers

import (
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nytopop/ssbd/models"
)

type AddVolumeForm struct {
	Name     string
	Backend  string // dir, sftp, aws, ceph, etc
	AuthUser string
	AuthPW   string
}

// GET /volumes/list
func VolumesList(db models.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		vols, err := db.GetVolumes()
		if err != nil {
			RenderErr(c, err)
			return
		}

		name := make([]byte, 64)
		for i := range name {
			name[i] = byte(rand.Intn(150))
		}

		err = db.InsertVolume(models.Volume{
			Name:    string(name),
			Backend: "one true backend",
		})
		if err != nil {
			RenderErr(c, err)
			return
		}

		c.HTML(http.StatusOK, "volumes/list.html", gin.H{
			"Volumes": vols,
		})
	}
}

// GET /volumes/add
func VolumesAdd() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "volumes/add.html", gin.H{})
	}
}

// POST /volumes/add
func TryVolumesAdd(db models.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		// validate the form
		// send a
	}
}
