package main

import (
	"log/slog"
	"net/http"
	"os"
	ws "web-chat/internal/delivery/websocket"
	"web-chat/internal/handler"
	"web-chat/internal/hub"
	"web-chat/internal/logger"
	"web-chat/internal/middleware"
	"web-chat/internal/service"
)

func main() {
	cfg := logger.Config{
		Level:  logger.LevelFromString(os.Getenv("LOG_LEVEL")),
		Format: logger.FormatFromString(os.Getenv("LOG_FORMAT")),
	}

	log := logger.NewLogger(cfg)
	slog.SetDefault(log)

	userService := service.NewUserMemory()
	userHandler := handler.NewUserHandler(userService)

	roomService := service.NewRoomMemory()
	roomHandler := handler.NewRoomHandler(roomService)

	messageService := service.NewMessageMemory()
	messageHandler := handler.NewMessageHandler(messageService)

	chatHub := hub.NewHub()
	go chatHub.Run()

	wsHandler := ws.ServeWS(chatHub, messageService, userService)

	mux := http.NewServeMux()

	// users
	mux.HandleFunc("POST /users", userHandler.Create)
	mux.HandleFunc("GET /users/{id}", userHandler.GetByID)
	mux.HandleFunc("PUT /users/{id}", userHandler.Update)
	mux.HandleFunc("DELETE /users/{id}", userHandler.Delete)

	// rooms
	mux.HandleFunc("POST /rooms", roomHandler.Create)
	mux.HandleFunc("GET /rooms/{id}", roomHandler.GetByID)
	mux.HandleFunc("GET /rooms", roomHandler.GetAll)
	mux.HandleFunc("PUT /rooms/{id}", roomHandler.Update)
	mux.HandleFunc("DELETE /rooms/{id}", roomHandler.Delete)

	// messages
	mux.HandleFunc("POST /rooms/{room_id}/messages", messageHandler.Create)
	mux.HandleFunc("GET /messages/{id}", messageHandler.GetByID)
	mux.HandleFunc("GET /rooms/{room_id}/messages", messageHandler.GetByRoomID)
	mux.HandleFunc("DELETE /messages/{id}", messageHandler.Delete)
	mux.HandleFunc("DELETE /rooms/{room_id}/messages", messageHandler.DeleteByRoom)

	// ws
	mux.HandleFunc("GET /ws/chat/{room_id}", wsHandler)

	slog.Info("server starting", "port", 8080)

	muxWithLogging := middleware.LoggingMiddleware(mux)
	if err := http.ListenAndServe(":8080", muxWithLogging); err != nil {
		slog.Error("server failed", "error", err)
	}

}
