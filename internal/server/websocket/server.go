package websocket

import (
	"net/http"
	"sync"

	"github.com/hoangminhphuc/goph-chat/common"
)

type WebSocketServer struct {
	name  string
	Rooms map[int]*Pool
	mu    sync.RWMutex
}

func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{
		name:  "websocket",
		Rooms: make(map[int]*Pool),
	}
}

func (s *WebSocketServer) Name() string {
	return s.name
}

func (s *WebSocketServer) InitFlags() {
	// No need for flags (for now)
}

func (s *WebSocketServer) Run() error {
	for _, room := range s.Rooms {
		go room.Start()
	}
	return nil
}

func (s *WebSocketServer) Stop() <-chan error {
	c := make(chan error, 1)
    go func() {
        s.mu.RLock()
        defer s.mu.RUnlock()
        for _, room := range s.Rooms {
            room.Stop()
        }
        c <- nil
        close(c)
    }()
    return c
}

func (s *WebSocketServer) GetRoom(roomID int) (*Pool, error) {
	s.mu.RLock()
	pool, ok := s.Rooms[roomID]
	s.mu.RUnlock()

	if !ok {
		return nil, common.ErrNotFound("room", roomID)
	} 
	
	return pool, nil
}

func (s *WebSocketServer) CreateRoom(roomID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.Rooms[roomID]; !exists {
		pool := NewPool(roomID)
		s.Rooms[roomID] = pool
	} else {
		return common.NewError("room already exists", http.StatusBadRequest)
	}
	
	return nil
}

