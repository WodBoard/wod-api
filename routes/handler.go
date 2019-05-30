package routes

import (
	"net/http"

	"github.com/WodBoard/wod-api/routes/storage"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Handler is defines the structure handling every route
type Handler struct {
	Storage *storage.Storage
	Lol     string
}

// AuthMiddleware serves as the middleware auth for every route but /login and /signup
func (h *Handler) AuthMiddleware(auths ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user needs to be signed in to access this service"})
			c.Abort()
			return
		}
		if len(auths) != 0 {
			authType := session.Get("authType")
			if authType == nil {
				c.JSON(http.StatusForbidden, gin.H{"error": "invalid request, restricted endpoint"})
				c.Abort()
				return
			}
		}
		// add session verification here, like checking if the user and authType
		// combination actually exists if necessary. Try adding caching this (redis)
		// since this middleware might be called a lot
		c.Next()
	}
}

// HandleRoutes defines the bindings between routes and functions
func (h *Handler) HandleRoutes(r *gin.Engine) {
	api := r.Group("/")
	{
		api.POST("/login", h.Login)
		api.POST("/signup", h.Signup)
	}
	apiAuth := api.Group("/")
	{
		apiAuth.Use(h.AuthMiddleware())
		apiAuth.GET("/logout", h.Logout)
		apiAuth.GET("/ping", h.Ping)
	}
}
