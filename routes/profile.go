package routes

import (
	"context"
	"log"
	"net/http"

	user "github.com/WodBoard/models/user/go"
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

// Profile is a basic endpoint just for example
func (h *Handler) Profile(c *gin.Context) {
	ctx := context.Background()
	claims := jwt.ExtractClaims(c)
	email, _ := claims[identityKey]

	user, err := h.Storage.GetUserByEmail(ctx, email.(string))
	if err != nil {
		log.Println(
			"err", err,
			"msg", "couldn't extract user informations",
			"email", email,
		)
		c.Status(http.StatusInternalServerError)
		return
	}

	err = h.OutputProtoMessage(c, user)
	if err != nil {
		return
	}

	c.Status(200)
}

// EditProfile is a basic endpoint just for example
func (h *Handler) EditProfile(c *gin.Context) {
	var req user.User
	ctx := context.Background()
	claims := jwt.ExtractClaims(c)
	email, _ := claims[identityKey]

	err := h.ParseProtoMessage(c, &req)
	if err != nil {
		return
	}

	err = h.Storage.UpdateUserByMail(ctx, email.(string), &req)
	if err != nil {
		log.Println(
			"err", err,
			"msg", "couldn't extract user informations",
			"email", email,
		)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(200)
}
