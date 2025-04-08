package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/meshapi/grpc-api-gateway-examples/gen"
)

type Update struct {
	UserID string
	Name   string
	Delete bool
}

type UserService struct {
	gen.UnimplementedUserServiceServer
	users   map[string]string
	streams map[string]chan Update
	mu      sync.Mutex
}

func NewUserService() *UserService {
	return &UserService{
		users:   make(map[string]string),
		streams: make(map[string]chan Update),
	}
}

func (s *UserService) AddUser(ctx context.Context, req *gen.AddUserRequest) (*gen.AddUserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.New().String()
	s.users[id] = req.Name
	go s.broadcastUpdate(Update{UserID: id, Name: req.Name, Delete: false})

	return &gen.AddUserResponse{Id: id}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *gen.DeleteUserRequest) (*gen.DeleteUserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	name, ok := s.users[req.Id]
	if !ok {
		return &gen.DeleteUserResponse{}, nil
	}

	delete(s.users, req.Id)
	go s.broadcastUpdate(Update{UserID: req.Id, Name: name, Delete: true})

	return &gen.DeleteUserResponse{}, nil
}

func (s *UserService) broadcastUpdate(update Update) {
	for _, updateChan := range s.streams {
		updateChan <- update
	}
}

func (s *UserService) UserStream(req *gen.UserStreamRequest, stream gen.UserService_UserStreamServer) error {
	streamID := uuid.New().String()
	updateChan := make(chan Update)

	s.mu.Lock()
	s.streams[streamID] = updateChan
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		delete(s.streams, streamID)
	}()

	for update := range updateChan {
		if update.Delete && !req.IncludeDeletions {
			continue
		}
		err := stream.Send(&gen.UserStreamResponse{
			Id:      update.UserID,
			Name:    update.Name,
			Deleted: update.Delete,
		})
		if err != nil {
			return fmt.Errorf("failed to send to peer: %w", err)
		}
	}

	return nil
}
