package service

import (
	"log"

	"github.com/Aysnine/unripe-bison/internal/types"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/gofiber/websocket/v2"
)

type client struct {
	nickname string
	// Add more data to this type if needed
}

type chatMessage struct {
	connection *websocket.Conn
	content    string
}

var clients = make(map[*websocket.Conn]client) // Note: although large maps with pointer-like types (e.g. strings) as keys are slow, using pointers themselves as keys is acceptable and fast
var register = make(chan *websocket.Conn)
var unregister = make(chan *websocket.Conn)
var chat = make(chan *chatMessage)

func runHub() {
	for {
		select {
		case connection := <-register:
			nickname := utils.ImmutableString(connection.Query("nickname"))
			if len(nickname) == 0 {
				nickname = "anonymous"
			}

			clients[connection] = client{nickname}

			log.Printf("[join] %v", nickname)

		case message := <-chat:
			nickname := clients[message.connection].nickname
			log.Printf("[message] %v: %v", nickname, message.content)

			// Send the message to all clients
			for connection := range clients {
				err := connection.WriteMessage(websocket.TextMessage, []byte(message.content))

				if err != nil {
					log.Println("write error:", err)
					connection.WriteMessage(websocket.CloseMessage, []byte{})
					connection.Close()
					delete(clients, connection)
				}
			}

		case connection := <-unregister:
			nickname := clients[connection].nickname

			// Remove the client from the hub
			delete(clients, connection)

			log.Printf("[leave] %v", nickname)
		}
	}
}

// Chat godoc
// @Summary global chat room
// @ID chat
// @Produce  websocket
// @Router /chat [get]
func SetupWebsocket_Chat(setup *types.SetupContext) {
	app := setup.App

	app.Use("/chat", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) { // Returns true if the client requested upgrade to the WebSocket protocol
			return c.Next()
		}
		return c.SendStatus(fiber.StatusUpgradeRequired)
	})

	go runHub()

	app.Get("/chat", websocket.New(func(connection *websocket.Conn) {
		// When the function returns, unregister the client and close the connection
		defer func() {
			unregister <- connection
			connection.Close()
		}()

		// Register the client
		register <- connection

		for {
			messageType, messageContent, err := connection.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println("read error:", err)
				}

				return // Calls the deferred function, i.e. closes the connection on error
			}

			if messageType == websocket.TextMessage {
				// Broadcast the received message
				chat <- &chatMessage{connection: connection, content: string(messageContent)}
			} else {
				log.Println("websocket message received of type", messageType)
			}
		}
	}))
}
