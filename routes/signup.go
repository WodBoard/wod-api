package routes

import (
	"context"
	"log"
	"net/http"

	user "github.com/WodBoard/models/user/go"
	"github.com/gin-gonic/gin"
)

// Signup is a basic endpoint just for example
func (h *Handler) Signup(c *gin.Context) {
	var req user.Signup
	ctx := context.Background()

	err := h.ParseProtoMessage(c, &req)
	if err != nil {
		return
	}

	login, err := h.Storage.GetLoginByEmail(ctx, req.GetEmail())
	if login != nil {
		log.Println(
			"msg", "user is already signed up",
			"email", req.GetEmail(),
		)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err = h.Storage.InsertLogin(ctx, &user.Login{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	})
	if err != nil {
		log.Println(
			"msg", "couldn't insert login in database",
			"err", err,
		)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = h.Storage.InsertUser(ctx, &user.User{
		Email:      req.GetEmail(),
		Firstname:  req.GetFirstname(),
		Lastname:   req.GetLastname(),
		Birthday:   req.GetBirthday(),
		PictureUrl: req.GetPictureUrl(),
		Height:     float64(req.GetHeight()),
		Weight:     float64(req.GetWeight()),
	})
	if err != nil {
		log.Println(
			"msg", "couldn't insert user in database",
			"err", err,
		)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}
