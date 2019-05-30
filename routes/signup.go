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

	err := c.ShouldBind(&req)
	if err != nil {
		log.Println("err", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	login, err := h.Storage.GetLoginByEmail(ctx, req.GetUsername())
	if login != nil {
		log.Println("msg", "user is already signed up", "email", req.GetUsername())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err = h.Storage.InsertLogin(ctx, &user.Login{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
	})
	if err != nil {
		log.Println("err", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = h.Storage.InsertUser(ctx, &user.User{
		Username:  req.GetUsername(),
		Firstname: req.GetFirstname(),
		Lastname:  req.GetLastname(),
	})
	if err != nil {
		log.Println("err", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}
