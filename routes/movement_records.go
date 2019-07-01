package routes

import (
	"context"
	"log"
	"net/http"

	movement_record "github.com/WodBoard/models/movement_record/go"
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

// ListMovementRecords is a GET endpoint that returns an array of movement records
// corresponding to the user's
func (h *Handler) ListMovementRecords(c *gin.Context) {
	ctx := context.Background()
	claims := jwt.ExtractClaims(c)
	email, _ := claims[identityKey]

	movements, err := h.Storage.ListMovementRecords(ctx, email.(string))
	if err != nil {
		log.Println(
			"err", err,
			"msg", "couldn't extract user's movement record",
			"email", email,
		)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(200, movements)
}

// AddMovementRecord is a POST endpoint that lets the user insert a
// movement record into our database
func (h *Handler) AddMovementRecord(c *gin.Context) {
	var req movement_record.MovementRecord
	ctx := context.Background()
	claims := jwt.ExtractClaims(c)
	email, _ := claims[identityKey]

	err := h.ParseProtoMessage(c, &req)
	if err != nil {
		return
	}

	t, err := h.Storage.GetMovementRecord(ctx, email.(string), req.GetName())
	if err == nil || t != nil {
		log.Println(
			"info", "can't insert the same training twice",
			"name", req.GetName(),
		)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err = h.Storage.InsertMovementRecord(ctx, email.(string), req)
	if err != nil {
		log.Println(
			"err", err,
			"msg", "couldn't insert movement record",
			"email", email,
			"name", req.GetName(),
		)
		c.Status(http.StatusInternalServerError)
		return
	}

	err = h.OutputProtoMessage(c, &req)
	if err != nil {
		return
	}

	c.Status(200)
}

// DeleteMovementRecord is a DELETE endpoint that lets the user delete a
// previously created movement record from our database
func (h *Handler) DeleteMovementRecord(c *gin.Context) {
	var req movement_record.DeleteMovementRecordRequest
	ctx := context.Background()
	claims := jwt.ExtractClaims(c)
	email, _ := claims[identityKey]

	err := h.ParseProtoMessage(c, &req)
	if err != nil {
		return
	}

	err = h.Storage.DeleteMovementRecord(ctx, email.(string), req.GetName())
	if err != nil {
		log.Println(
			"err", err,
			"msg", "couldn't delete user movement record",
			"email", email,
			"name", req.GetName(),
		)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(200)
}
