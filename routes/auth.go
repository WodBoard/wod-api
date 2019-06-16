package routes

import (
	"context"
	"log"

	user "github.com/WodBoard/models/user/go"
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

// Authenticator is part of the middleware auth
func (h *Handler) Authenticator(c *gin.Context) (interface{}, error) {
	ctx := context.Background()
	var loginReq user.Login
	if err := c.ShouldBind(&loginReq); err != nil {
		return "", jwt.ErrMissingLoginValues
	}
	email := loginReq.GetEmail()
	password := loginReq.GetPassword()

	login, err := h.Storage.GetLoginByEmail(ctx, email)
	if err != nil {
		log.Println(
			"err", err,
			"msg", "Failed to get login",
			"email", email,
		)
		return nil, jwt.ErrFailedAuthentication
	}
	if password != login.Password {
		log.Println("msg", "invalid password")
		return nil, jwt.ErrFailedAuthentication
	}

	user, err := h.Storage.GetUserByEmail(ctx, email)
	if err != nil {
		log.Println(
			"err", err,
			"msg", "Failed to get user",
			"email", email,
		)
		return nil, jwt.ErrFailedAuthentication
	}
	return user, nil
}

// Authorizator is a middleware auth func
func (h *Handler) Authorizator(data interface{}, c *gin.Context) bool {
	if v, ok := data.(*user.User); ok {
		ctx := context.Background()
		_, err := h.Storage.GetUserByEmail(ctx, v.GetEmail())
		return err == nil
	}
	return false
}

// Unauthorized is the default func when access is denied
func (h *Handler) Unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}
