package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

const (
	host    = "localhost:8080"
	token   = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxOSwiZXhwIjoxNzg0NTQyMDQ5fQ.1_aSFCwgruLY_1qvprPeiQteUK2kexvQtBlMe8MG1ws"
	room_id = "93"
)

func connectClient(id int) (*websocket.Conn, error) {
	u := url.URL{
		Scheme:   "ws",
		Host:     host,
		Path:     "/ws/chat/" + room_id,
		RawQuery: "token=" + token,
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("client %d diel %w", id, err)
	}

	return conn, nil
}

func runClient(id int, wg *sync.WaitGroup, msgCounter *int64, connectTimes *[]time.Duration, errCounter *int64, mu *sync.Mutex) {
	defer wg.Done()

	connectStart := time.Now()
	conn, err := connectClient(id)
	connectDuration := time.Since(connectStart)

	if err != nil {
		log.Printf("client %d failed to connect: %v", id, err)
		atomic.AddInt64(errCounter, 1)
		return
	}
	defer conn.Close()

	mu.Lock()
	*connectTimes = append(*connectTimes, connectDuration)
	mu.Unlock()

	for i := range 100 {
		msg := map[string]string{
			"type": "content",
			"text": fmt.Sprintf("hello from client %d, msg %d", id, i),
		}

		data, _ := json.Marshal(msg)

		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("client %d write error: %v", id, err)
			atomic.AddInt64(errCounter, 1)
			return
		}

		atomic.AddInt64(msgCounter, 1)
		time.Sleep(100 * time.Millisecond)
	}
}

func percentile(sorted []time.Duration, p float64) time.Duration {
	if len(sorted) == 0 {
		return 0
	}

	idx := int(float64(len(sorted)) * p)
	if idx >= len(sorted) {
		idx = len(sorted) - 1
	}

	return sorted[idx]
}

func main() {
	numClients := 10_000

	var wg sync.WaitGroup
	var msgCount int64
	var errCount int64
	var mu sync.Mutex
	connectTimes := make([]time.Duration, 0, numClients)

	start := time.Now()

	for i := range numClients {
		wg.Add(1)
		go runClient(i, &wg, &msgCount, &connectTimes, &errCount, &mu)
	}

	wg.Wait()

	elapsed := time.Since(start)

	slices.Sort(connectTimes)

	fmt.Println("---")
	fmt.Printf("clients: %d\n", numClients)
	fmt.Printf("messages sent: %d\n", msgCount)
	fmt.Printf("errors: %d\n", errCount)
	fmt.Printf("elapsed: %s\n", elapsed)
	fmt.Printf("msg/sec: %.2f\n", float64(msgCount)/elapsed.Seconds())

	fmt.Println("--- connect latency ---")
	fmt.Printf("min: %s\n", connectTimes[0])
	fmt.Printf("p50: %s\n", percentile(connectTimes, 0.50))
	fmt.Printf("p95: %s\n", percentile(connectTimes, 0.95))
	fmt.Printf("p99: %s\n", percentile(connectTimes, 0.99))
	fmt.Printf("max: %s\n", connectTimes[len(connectTimes)-1])
}
