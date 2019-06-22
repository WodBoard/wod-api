package routes

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"

	training "github.com/WodBoard/models/training/go"
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	jsonpb "github.com/golang/protobuf/jsonpb"
)

// ListTrainings is a GET endpoint that returns an array of trainings
// corresponding to the user's
func (h *Handler) ListTrainings(c *gin.Context) {
	ctx := context.Background()
	claims := jwt.ExtractClaims(c)
	email, _ := claims[identityKey]

	trainings, err := h.Storage.ListTrainings(ctx, email.(string))
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

// AddTraining is a PUT endpoint that lets the user insert a
// personal training into our database
func (h *Handler) AddTraining(c *gin.Context) {
	var req training.Training
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
			"msg", "couldn't unmarshal add training request",
		)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	t, err := h.Storage.GetTraining(ctx, email.(string), req.GetName())
	if err == nil || t != nil {
		log.Println(
			"info", "can't insert the same training twice",
			"name", req.GetName(),
		)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err = h.Storage.InsertTraining(ctx, email.(string), req)
	if err != nil {
		log.Println(
			"err", err,
			"msg", "couldn't insert user training",
			"email", email,
			"name", req.GetName(),
		)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(200, &req)
}

// EditTraining is a POST endpoint that lets the user edit
// a previously created training
func (h *Handler) EditTraining(c *gin.Context) {
	var req training.Training
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
			"msg", "couldn't unmarshal add training request",
		)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err = h.Storage.UpdateTraining(ctx, email.(string), req.GetName(), req)
	if err != nil {
		log.Println(
			"err", err,
			"msg", "couldn't edit user training",
			"email", email,
			"name", req.GetName(),
		)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(200, &req)
}

// DeleteTraining is a DELETE endpoint that lets the user delete a
// previously created training from our database
func (h *Handler) DeleteTraining(c *gin.Context) {
	var req training.DeleteTrainingRequest
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
			"msg", "couldn't unmarshal add training request",
		)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err = h.Storage.DeleteTraining(ctx, email.(string), req.GetName())
	if err != nil {
		log.Println(
			"err", err,
			"msg", "couldn't delete user training",
			"email", email,
			"name", req.GetName(),
		)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(200)
}
