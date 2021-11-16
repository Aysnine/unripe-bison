package service

import (
	"log"

	"github.com/Aysnine/unripe-bison/internal/types"
	"github.com/gofiber/websocket/v2"
)

// Chat godoc
// @Summary global chat room
// @ID chat
// @Produce  websocket
// @Router /chat [get]
func SetupWebsocket_Chat(setup *types.SetupContext) {
	app := setup.App

	app.Get("/chat", websocket.New(func(c *websocket.Conn) {
		// Websocket logic
		for {
			messageType, msg, err := c.ReadMessage()
			if err != nil {
				break
			}
			log.Printf("Read: %s", msg)

			err = c.WriteMessage(messageType, msg)
			if err != nil {
				break
			}
			log.Println("Error:", err)
		}
	}))
}
