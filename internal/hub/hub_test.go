package hub

import (
	"sync"
	"testing"
	"time"
)

type mockOnlineRepository struct {
	mtx   sync.Mutex
	rooms map[int]map[int]struct{}
}

func (m *mockOnlineRepository) AddOnline(roomID, userID int) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if m.rooms[roomID] == nil {
		m.rooms[roomID] = make(map[int]struct{})
	}

	m.rooms[roomID][userID] = struct{}{}

	return nil
}

func (m *mockOnlineRepository) RemoveOnline(roomID, userID int) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	delete(m.rooms[roomID], userID)

	return nil
}

func (m *mockOnlineRepository) GetOnlineCount(roomID int) (int, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	return len(m.rooms[roomID]), nil
}

func (m *mockOnlineRepository) GetOnlineUsers(roomID int) ([]int, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	ids := make([]int, 0, len(m.rooms[roomID]))

	for id := range m.rooms[roomID] {
		ids = append(ids, id)
	}

	return ids, nil
}

func (m *mockOnlineRepository) IsOnline(roomID, userID int) (bool, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if _, ok := m.rooms[roomID][userID]; ok {
		return true, nil
	}

	return false, nil
}

func newTestHub() *Hub {
	h := NewHub(&mockOnlineRepository{
		rooms: make(map[int]map[int]struct{}),
	})
	go h.Run()
	time.Sleep(5 * time.Millisecond)
	return h
}

func newClient(userID, roomID int) *Client {
	return &Client{
		UserID: userID,
		RoomID: roomID,
		Send:   make(chan []byte, 1),
	}
}

func newFullSendClient(userID, roomID int) *Client {
	return &Client{
		UserID: userID,
		RoomID: roomID,
		Send:   make(chan []byte),
	}
}

func TestHub_Register(t *testing.T) {
	h := newTestHub()
	c := newClient(1, 1)

	h.Register <- c
	time.Sleep(10 * time.Millisecond)

	count := h.GetOnlineCount(1)
	if count != 1 {
		t.Errorf("expected 1, got %d", count)
	}
}

func TestHub_Unregister(t *testing.T) {
	h := newTestHub()
	c := newClient(1, 1)

	h.Register <- c
	time.Sleep(20 * time.Millisecond)

	count1 := h.GetOnlineCount(1)
	if count1 != 1 {
		t.Errorf("before unregister: expected 1, got %d", count1)
	}

	h.Unregister <- c
	time.Sleep(20 * time.Millisecond)

	count2 := h.GetOnlineCount(1)
	if count2 != 0 {
		t.Errorf("after unregister: expected 0, got %d", count2)
	}
}

func TestHub_Unregister_NotExists(t *testing.T) {
	h := newTestHub()

	c := newClient(1, 1)
	h.Unregister <- c
	time.Sleep(20 * time.Millisecond)
}

func TestHub_Broadcast_SameRoom(t *testing.T) {
	h := newTestHub()
	c1 := newClient(1, 1)
	c2 := newClient(2, 1)

	h.Register <- c1
	h.Register <- c2
	time.Sleep(10 * time.Millisecond)

	h.Broadcast <- BroadcastMsg{
		RoomID: 1,
		Data:   []byte("hello"),
	}

	select {
	case msg := <-c1.Send:
		if string(msg) != "hello" {
			t.Errorf("c1: expected `hello`, got `%s`", msg)
		}
	case <-time.After(time.Second):
		t.Error("c1: timeout waiting for message")
	}

	select {
	case msg := <-c2.Send:
		if string(msg) != "hello" {
			t.Errorf("c2: expected `hello`, got `%s`", msg)
		}
	case <-time.After(time.Second):
		t.Error("c2: timeout waiting for message")
	}
}

func TestHub_Broadcast_DiffRoom(t *testing.T) {
	h := newTestHub()

	c1 := newClient(1, 1)
	c2 := newClient(2, 2)

	h.Register <- c1
	h.Register <- c2
	time.Sleep(10 * time.Millisecond)

	h.Broadcast <- BroadcastMsg{
		RoomID: 1,
		Data:   []byte("hello"),
	}

	select {
	case msg := <-c1.Send:
		if string(msg) != "hello" {
			t.Errorf("c1: expected `hello`, got `%s`", msg)
		}
	case <-time.After(time.Second):
		t.Error("c1: timeout waiting for message")
	}

	select {
	case <-c2.Send:
		t.Error("c2 shouldn't receive message")
	default:
	}
}

func TestHub_GetOnlineCount(t *testing.T) {
	h := newTestHub()

	c1 := newClient(1, 1)
	c2 := newClient(2, 1)
	c3 := newClient(3, 2)

	h.Register <- c1
	h.Register <- c2
	h.Register <- c3
	time.Sleep(10 * time.Millisecond)

	count1 := h.GetOnlineCount(1)
	if count1 != 2 {
		t.Errorf("room 1: expected 2, got %d", count1)
	}

	count2 := h.GetOnlineCount(2)
	if count2 != 1 {
		t.Errorf("room 2: expected 1, got %d", count2)
	}

	count3 := h.GetOnlineCount(999)
	if count3 != 0 {
		t.Errorf("room 999: expected 0, got %d", count3)
	}
}

func TestHub_CheckOnline(t *testing.T) {
	h := newTestHub()

	c := newClient(1, 1)
	h.Register <- c
	time.Sleep(10 * time.Millisecond)

	exists1 := h.CheckOnline(1, 1)
	if !exists1 {
		t.Errorf("room 1, user 1: expected: true, got false")
	}

	exists2 := h.CheckOnline(2, 2)
	if exists2 {
		t.Errorf("room 2, user 2: expected: false, got true")
	}

	exists3 := h.CheckOnline(1000, 10)
	if exists3 {
		t.Errorf("room 10, user 1000: expected: false, got true")
	}
}

func TestHub_GetUsersInRoom(t *testing.T) {
	h := newTestHub()

	c1 := newClient(1, 1)
	c2 := newClient(2, 1)
	c3 := newClient(3, 1)

	c4 := newClient(4, 2)
	c5 := newClient(5, 2)
	c6 := newClient(6, 3)

	h.Register <- c1
	h.Register <- c2
	h.Register <- c3
	h.Register <- c4
	h.Register <- c5
	h.Register <- c6
	time.Sleep(10 * time.Millisecond)

	clients1 := h.GetUsersInRoom(1)
	if len(clients1) != 3 {
		t.Errorf("room 1: expected 3, got %d", len(clients1))
	}

	clients2 := h.GetUsersInRoom(2)
	if len(clients2) != 2 {
		t.Errorf("room 2: expected 2, got %d", len(clients2))
	}

	clients3 := h.GetUsersInRoom(3)
	if len(clients3) != 1 {
		t.Errorf("room 3: expected 1, got %d", len(clients3))
	}
}

func TestHub_Broadcast_FullSend(t *testing.T) {
	h := newTestHub()
	c := newFullSendClient(1, 1)

	h.Register <- c
	time.Sleep(10 * time.Millisecond)

	h.Broadcast <- BroadcastMsg{
		RoomID: 1,
		Data:   []byte("hello"),
	}

	count := h.GetOnlineCount(1)
	if count != 0 {
		t.Errorf("expected 0, got %d", count)
	}
}

func TestHub_Concurrent_Register(t *testing.T) {
	h := newTestHub()

	var wg sync.WaitGroup
	for i := range 100 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			c := newClient(id, 1)
			h.Register <- c
		}(i)
	}
	wg.Wait()

	count := h.GetOnlineCount(1)
	if count != 100 {
		t.Errorf("expected 100, got %d", count)
	}
}
