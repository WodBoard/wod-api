package routes

import "github.com/gin-gonic/gin"

// Signup is a basic endpoint just for example
func (h *Handler) Signup(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": h.Lol,
	})
}
