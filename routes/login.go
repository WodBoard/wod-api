package routes

import (
	"github.com/gin-gonic/gin"
)

// Login is a basic endpoint just for example
func (h *Handler) Login(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": h.Lol,
	})
}

// Logout is a basic endpoint just for example
func (h *Handler) Logout(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": h.Lol,
	})
}
