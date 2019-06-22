package routes

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"

	user "github.com/WodBoard/models/user/go"
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	jsonpb "github.com/golang/protobuf/jsonpb"
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
	c.JSON(200, user)
}

// EditProfile is a basic endpoint just for example
func (h *Handler) EditProfile(c *gin.Context) {
	var req user.User
	ctx := context.Background()
	claims := jwt.ExtractClaims(c)
	email, _ := claims[identityKey]

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(
			"err", err,
			"msg", "couldn't read body of the request",
		)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	err = jsonpb.UnmarshalString(string(body), &req)
	if err != nil {
		log.Println(
			"err", err,
			"msg", "couldn't unmarshal signup request",
		)
		c.AbortWithStatus(http.StatusBadRequest)
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
