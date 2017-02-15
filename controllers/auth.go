package controllers

import (
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AuthCheckpoint checks cookies to see if user is signed in. If they are,
// it sets the logged in user in the current request's context chain.
func AuthCheckpoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		// cookies
		cookies := sessions.Default(c)

		// check if sign in cookie exists
		if cookies.Get("userid") != nil {
			c.Set("userid", cookies.Get("userid"))
			c.Set("name", cookies.Get("name"))
		}
	}
}

// AuthSignedIn checks if a user is signed in, if not it will
// redirect to /auth/sign-in.
func AuthSignedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		// if userid exists in the current context, we are already signed in
		_, exists := c.Get("userid")
		if !exists {
			c.Redirect(302, "/auth/sign-in")
		}
	}
}

// AuthSignIn serves the sign in page HTML.
// GET /auth/sign-in
func AuthSignIn(c *gin.Context) {
	c.HTML(http.StatusOK, "main/main.html", gin.H{})
}

// TryAuthSignIn attempts to sign in a user by reading data from the POST
// form and setting session cookies corresponding to a user in the database.
// POST /auth/sign-in
func AuthTrySignIn(c *gin.Context) {
}

// TryAuthSignOut attempts to sign out the current user by clearing any
// session cookies, and deleting the session itself.
// GET /auth/sign-out
func AuthTrySignOut(c *gin.Context) {
}
