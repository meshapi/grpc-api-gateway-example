package main

import (
	"io"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/meshapi/grpc-api-gateway-examples/gen"
)

type ChatService struct {
	gen.UnimplementedChatServiceServer
	users map[string]gen.ChatService_ChatServer
	lock  sync.Mutex
}

func NewChatService() *ChatService {
	return &ChatService{
		users: make(map[string]gen.ChatService_ChatServer),
	}
}

func (c *ChatService) Chat(stream gen.ChatService_ChatServer) error {
	clientID := uuid.New().String()

	// Register the user
	c.lock.Lock()
	c.users[clientID] = stream
	c.lock.Unlock()

	defer func() {
		// Unregister the user when the stream is closed
		c.lock.Lock()
		delete(c.users, clientID)
		c.lock.Unlock()
	}()

	for {
		// Receive a message from the client
		req, err := stream.Recv()
		if err == io.EOF {
			// Client closed the stream
			return nil
		}
		if err != nil {
			return err
		}

		// Broadcast the message to all connected users
		c.lock.Lock()
		for addr, userStream := range c.users {
			if addr == clientID {
				continue // Don't send the message back to the sender
			}
			err := userStream.Send(&gen.ChatResponse{
				User: req.Name,
				Text: req.Text,
			})
			if err != nil {
				// Handle the error (e.g., log it, remove the user from the map, etc.)
				log.Printf("Error sending message to %s: %v\n", addr, err)
			}
		}
		c.lock.Unlock()
	}
}
