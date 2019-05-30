package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/WodBoard/wod-api/storage"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Handler is defines the structure handling every route
type Handler struct {
	Storage *storage.Storage
	Lol     string
}

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

// Ping is a basic endpoint just for example
func (h *Handler) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": h.Lol,
	})
}

// Ping is a basic endpoint just for example
func (h *Handler) Login(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": h.Lol,
	})
}

// Ping is a basic endpoint just for example
func (h *Handler) Logout(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": h.Lol,
	})
}

// Ping is a basic endpoint just for example
func (h *Handler) Signup(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": h.Lol,
	})
}

// HandleRoutes defines the bindings between routes and functions
func (h *Handler) HandleRoutes(r *gin.Engine) {
	api := r.Group("/")
	{
		api.POST("/login", h.Login)
		api.POST("/signup", h.Signup)
	}
	apiAuth := api.Group("/")
	apiAuth.Use(h.AuthMiddleware())
	apiAuth.GET("/logout", h.Logout)
	apiAuth.GET("/ping", h.Ping)
}

func main() {
	ctx := context.Background()

	mongoClient, err := mongo.NewClient(options.Client().
		ApplyURI(os.Getenv("MONGO_URI")).
		SetAuth(options.Credential{
			Username: os.Getenv("MONGO_USERNAME"),
			Password: os.Getenv("MONGO_PASSWORD"),
		}),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't create mongodb instance", err)
		os.Exit(2)
	}

	err = mongoClient.Connect(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't connect to mongodb instance", err)
		os.Exit(2)
	}

	storage := storage.NewStorage(mongoClient)
	handler := &Handler{
		Lol:     "lolilol",
		Storage: storage,
	}

	r := gin.Default()
	handler.HandleRoutes(r)

	r.Run(os.Getenv("LISTEN_ADDR"))
}
