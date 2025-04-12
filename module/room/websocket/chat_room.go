package websocket

import (
	"sync"
)

type RoomCenter struct {
	Rooms map[string]*Pool
	mu    sync.RWMutex
}

func NewRoomCenter() *RoomCenter {
	return &RoomCenter{
		Rooms: make(map[string]*Pool),
	}
}

var roomManager = NewRoomCenter()

func GetRoomCenter() *RoomCenter {
	return roomManager
}


func (rm *RoomCenter) GetRoom(roomID string) *Pool {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	pool, ok := rm.Rooms[roomID]

	if !ok {
		pool = NewPool(roomID)
		rm.Rooms[roomID] = pool
		go pool.Start()
	}
	return pool
}