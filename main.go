package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

func main() {
	server := gin.Default()
	server.LoadHTMLGlob("template/html/*")
	server.Static("/assets", "./template/assets")
	server.GET("/", DefaultPage)
	server.GET("/login", LoginPage)
	server.POST("/login", LoginAuth)
	server.GET("/chat", ChatPage)

	handler := melody.New()
	server.GET("/ws", func(c *gin.Context) {
		handler.HandleRequest(c.Writer, c.Request)
		fmt.Println(c.Writer, c.Request)
	})

	handler.HandleMessage(func(s *melody.Session, msg []byte) {
		handler.Broadcast(msg)
		fmt.Println("msg:", msg)
	})

	handler.HandleConnect(func(session *melody.Session) {
		id := session.Request.URL.Query().Get("id")
		handler.Broadcast((NewMessage("other", id, "join the chat").GetByteMessage()))
		fmt.Println("id:", id)
	})

	handler.HandleClose(func(session *melody.Session, i int, s string) error {
		id := session.Request.URL.Query().Get("id")
		handler.Broadcast(NewMessage("other", id, "exit the chat").GetByteMessage())
		fmt.Println("id:", id)
		return nil
	})
	server.Run(":8080")
}
