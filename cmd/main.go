package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	ws "web-chat/internal/delivery/websocket"
	"web-chat/internal/handler"
	"web-chat/internal/hub"
	"web-chat/internal/logger"
	"web-chat/internal/middleware"
	"web-chat/internal/repository"
	"web-chat/internal/service"
)

func main() {
	cfg := logger.Config{
		Level:  logger.LevelFromString(os.Getenv("LOG_LEVEL")),
		Format: logger.FormatFromString(os.Getenv("LOG_FORMAT")),
	}

	log := logger.NewLogger(cfg)
	slog.SetDefault(log)

	db := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s?sslmode=disable", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))

	pool, err := repository.NewPool(context.Background(), db)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	userRepo := repository.NewUserPG(pool)
	roomRepo := repository.NewRoomPG(pool)
	msgRepo := repository.NewMessagePG(pool)

	userService := service.NewUserMemory(userRepo)
	roomService := service.NewRoomMemory(roomRepo)
	messageService := service.NewMessageMemory(msgRepo)

	userHandler := handler.NewUserHandler(userService)
	roomHandler := handler.NewRoomHandler(roomService)
	messageHandler := handler.NewMessageHandler(messageService)

	chatHub := hub.NewHub()
	go chatHub.Run()

	wsHandler := ws.ServeWS(chatHub, messageService, userService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /users", userHandler.Create)
	mux.HandleFunc("GET /users/{id}", userHandler.GetByID)
	mux.HandleFunc("PUT /users/{id}", userHandler.Update)
	mux.HandleFunc("DELETE /users/{id}", userHandler.Delete)

	mux.HandleFunc("POST /rooms", roomHandler.Create)
	mux.HandleFunc("GET /rooms/{id}", roomHandler.GetByID)
	mux.HandleFunc("GET /rooms", roomHandler.GetAll)
	mux.HandleFunc("PUT /rooms/{id}", roomHandler.Update)
	mux.HandleFunc("DELETE /rooms/{id}", roomHandler.Delete)

	mux.HandleFunc("POST /rooms/{room_id}/messages", messageHandler.Create)
	mux.HandleFunc("GET /messages/{id}", messageHandler.GetByID)
	mux.HandleFunc("GET /rooms/{room_id}/messages", messageHandler.GetByRoomID)
	mux.HandleFunc("DELETE /messages/{id}", messageHandler.Delete)
	mux.HandleFunc("DELETE /rooms/{room_id}/messages", messageHandler.DeleteByRoom)

	mux.HandleFunc("GET /ws/chat/{room_id}", wsHandler)

	slog.Info("server starting", "port", 8080)

	muxWithLogging := middleware.LoggingMiddleware(mux)
	if err := http.ListenAndServe(":8080", muxWithLogging); err != nil {
		slog.Error("server failed", "error", err)
	}

}
