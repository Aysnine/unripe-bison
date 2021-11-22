package service

import (
	"context"
	"log"

	"github.com/Aysnine/unripe-bison/internal/types"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/gofiber/websocket/v2"
)

type Client struct {
	nickname string
}

type IntoMessage struct {
	connection *websocket.Conn
	content    string
}

type SendMessage struct {
	channel string
	content string
}

var clients = make(map[*websocket.Conn]Client) // Note: although large maps with pointer-like types (e.g. strings) as keys are slow, using pointers themselves as keys is acceptable and fast
var register = make(chan *websocket.Conn)
var unregister = make(chan *websocket.Conn)
var into = make(chan *IntoMessage)
var send = make(chan *SendMessage)

func runHub(setupContext *types.SetupContext) {
	chatRedis := setupContext.ChatRedis

	// There is no error because go-redis automatically reconnects on error.
	pubsub := chatRedis.Subscribe(context.Background(), "channel:global")
	// Close the subscription when we are done.
	defer pubsub.Close()

	// Receive pub messages
	go func() {
		for {
			msg, err := pubsub.ReceiveMessage(context.Background())
			if err != nil {
				log.Println("receive redis pub error:", err)
			}
			send <- &SendMessage{channel: msg.Channel, content: msg.Payload}
		}
	}()

	for {
		select {
		case connection := <-register:
			nickname := utils.ImmutableString(connection.Query("nickname"))
			if len(nickname) == 0 {
				nickname = "anonymous"
			}

			client := Client{nickname}

			clients[connection] = client

			log.Printf("[join] %v", nickname)

			err := chatRedis.Publish(context.Background(), "channel:global", string("["+client.nickname+"] 🔵")).Err()
			if err != nil {
				log.Println("redis publish error:", err)
			}

		case message := <-into:
			client := clients[message.connection]

			log.Printf("[into message] %v: %v", client.nickname, message.content)

			err := chatRedis.Publish(context.Background(), "channel:global", string("["+client.nickname+"] 💬 "+message.content)).Err()
			if err != nil {
				log.Println("redis publish error:", err)
			}

		case message := <-send:
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
			client := clients[connection]

			// Remove the client from the hub
			delete(clients, connection)

			log.Printf("[leave] %v", client.nickname)

			err := chatRedis.Publish(context.Background(), "channel:global", string("["+client.nickname+"] 🔴")).Err()
			if err != nil {
				log.Println("redis publish error:", err)
			}
		}
	}
}

// Chat godoc
// @Summary global chat room
// @ID chat
// @Produce  json
// @Router /chat [get]
func SetupWebsocket_Chat(setupContext *types.SetupContext) {
	app := setupContext.App

	app.Use("/chat", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) { // Returns true if the client requested upgrade to the WebSocket protocol
			return c.Next()
		}
		return c.SendStatus(fiber.StatusUpgradeRequired)
	})

	// TODO remove condition when test mocked
	if setupContext.ChatRedis != nil {
		go runHub(setupContext)
	}

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
				into <- &IntoMessage{connection: connection, content: string(messageContent)}
			} else {
				log.Println("websocket message received of type", messageType)
			}
		}
	}))
}
