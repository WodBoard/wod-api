package routes

import "github.com/gin-gonic/gin"

// Ping is a basic endpoint just for example
func (h *Handler) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": h.Lol,
	})
}
