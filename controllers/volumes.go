package controllers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
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

		name := make([]byte, 32)
		rand.Read(name)
		sum := sha256.Sum256(name)
		hash := base64.StdEncoding.EncodeToString(sum[:])

		err = db.InsertVolume(models.Volume{
			Name:    hash,
			Backend: 0,
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
