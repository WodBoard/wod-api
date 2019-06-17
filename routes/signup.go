package routes

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"

	user "github.com/WodBoard/models/user/go"
	"github.com/gin-gonic/gin"
	jsonpb "github.com/golang/protobuf/jsonpb"
)

// Signup is a basic endpoint just for example
func (h *Handler) Signup(c *gin.Context) {
	var req user.Signup
	ctx := context.Background()

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
		log.Println("err", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = h.Storage.InsertUser(ctx, &user.User{
		Email:      req.GetEmail(),
		Username:   req.GetUsername(),
		Firstname:  req.GetFirstname(),
		Lastname:   req.GetLastname(),
		Birthday:   req.GetBirthday(),
		PictureUrl: req.GetPictureUrl(),
		Height:     float64(req.GetHeight()),
		Weight:     float64(req.GetWeight()),
	})
	if err != nil {
		log.Println("err", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}
