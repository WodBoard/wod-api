package routes

import (
	"context"
	"log"
	"net/http"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

// Trainings is a GET endpoint that returns an array of trainings
// corresponding to the user's
func (h *Handler) Trainings(c *gin.Context) {
	ctx := context.Background()
	claims := jwt.ExtractClaims(c)
	email, _ := claims[identityKey]

	trainings, err := h.Storage.GetTrainingsByEmail(ctx, email.(string))
	if err != nil {
		log.Println(
			"err", err,
			"msg", "couldn't extract user's trainings",
			"email", email,
		)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(200, trainings)
}
