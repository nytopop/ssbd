// Package web provides web server functionality for ssbd.
package web

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/nytopop/ssbd/config"
	"github.com/nytopop/ssbd/controllers"
	"github.com/nytopop/ssbd/models"
)

// StartServer starts pouring gin.
func StartServer(db models.Handler) error {
	//	gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)
	pub := gin.New()

	// Load HTML templates
	pub.LoadHTMLGlob(config.CFG.Srv.ResourceDir + "/templates/**/*.html")

	// Generate session auth and encryption keys
	// Restarting the server will INVALIDATE existing sessions
	rnd := make([]byte, 16)
	_, err := rand.Read(rnd)
	if err != nil {
		return err
	}

	auth := sha512.Sum512(rnd[:8])
	enc := sha256.Sum256(rnd[8:16])
	store := sessions.NewCookieStore(auth[:], enc[:])

	// Misc Middleware
	pub.Use(controllers.Logger())
	pub.Use(gin.Recovery())
	pub.Use(sessions.Sessions("ssbd", store))

	// Auth Middleware
	pub.Use(controllers.AuthCheckpoint())
	//users := pub.Group("/", controllers.AuthSignedIn())

	// Routes
	pub.Static("/static", config.CFG.Srv.ResourceDir+"/static")

	pub.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/dash/overview")
	})

	pub.GET("/auth/sign-in", controllers.AuthSignIn())
	pub.POST("/auth/sign-in", controllers.AuthTrySignIn(db))
	pub.GET("/auth/sign-out", controllers.AuthTrySignOut())

	// TODO : switch the following to the users group

	// Dashboard
	pub.GET("/dash/overview", controllers.DashOverview(db))

	// Add, delete, modify servers
	//users.GET("/servers", controllers.Servers)
	//users.GET("/servers/add", controllers.ServersAdd)
	//users.POST("/servers/add", controllers.ServersTryAdd)
	//users.GET("/servers/del/:serverid", controllers.ServersTryDel)

	// Add, delete, modify backup jobs
	//users.GET("/jobs", controllers.Jobs)
	//users.GET("/jobs/queue", controllers.JobsQueue)
	//users.GET("/jobs/add", controllers.JobsAdd)
	//users.POST("/jobs/add", controllers.JobsTryAdd)
	//users.GET("/jobs/del/:jobid", controllers.JobsTryDel)

	// View job history
	//users.GET("/history/:page", controllers.History)
	//users.GET("/history/:jobid", controllers.HistoryJobID)

	// Browse backups
	//users.GET("/browse/:jobid", controllers.BrowseJobID)

	// Server administration
	//users.GET("/admin", controllers.Admin)
	//users.GET("/admin/users", controllers.AdminUsers)
	//users.POST("/admin/users/add", controllers.AdminTryUsersAdd)
	//users.GET("/admin/users/del/:userid", controllers.AdminTryUsersDel)

	pub.Run(config.CFG.Srv.Listen)

	return nil
}
