package routes

import (
	"log"
	"time"

	user "github.com/WodBoard/models/user/go"
	"github.com/WodBoard/wod-api/routes/storage"
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	jsonpb "github.com/golang/protobuf/jsonpb"
)

var identityKey = "id"

// Handler is defines the structure handling every route
type Handler struct {
	Storage   *storage.Storage
	engine    *gin.Engine
	marshaler jsonpb.Marshaler
	addr      string
}

// NewHandler returns a fresh instance of Handler struct
func NewHandler(storage *storage.Storage, addr string) *Handler {
	e := gin.Default()
	return &Handler{
		Storage: storage,
		engine:  e,
		addr:    addr,
		marshaler: jsonpb.Marshaler{
			EnumsAsInts: true,
			EmitDefaults: true,
		},
	}
}

// HandleRoutes defines the bindings between routes and functions
func (h *Handler) HandleRoutes() {
	// JWT middleware
	secretKey := "MyRandomHashingKeyForWod9876543210"
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "wod",
		Key:         []byte(secretKey),
		Timeout:     time.Hour * 24,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*user.User); ok {
				return jwt.MapClaims{
					identityKey: v.GetEmail(),
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &user.User{
				Email: claims[identityKey].(string),
			}
		},
		Authenticator: h.Authenticator,
		Authorizator:  h.Authorizator,
		Unauthorized:  h.Unauthorized,
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
	if err != nil {
		log.Fatalf("Cannot initialize middleware: %s\n", err.Error())
	}

	api := h.engine.Group("/")
	{
		api.POST("/login", authMiddleware.LoginHandler)
		api.POST("/signup", h.Signup)
	}
	auth := api.Group("/")
	{
		auth.Use(authMiddleware.MiddlewareFunc())
		auth.GET("/profile", h.Profile)
		auth.PUT("/profile", h.EditProfile)
		auth.GET("/trainings", h.ListTrainings)
		auth.PUT("/trainings", h.AddTraining)
		auth.POST("/trainings", h.EditTraining)
		auth.DELETE("/trainings", h.DeleteTraining)
	}
	h.engine.Run(h.addr)
}
