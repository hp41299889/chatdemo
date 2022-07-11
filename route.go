package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func DefaultPage(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "wtf")
}

func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func ChatPage(c *gin.Context) {
	nickname := c.Query("nickname")
	c.HTML(http.StatusOK, "chat.html", gin.H{
		"nickname": nickname,
	})
}
