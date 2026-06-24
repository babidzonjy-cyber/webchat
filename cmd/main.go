package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
	"web-chat/internal/auth"
	ws "web-chat/internal/delivery/websocket"
	"web-chat/internal/handler"
	"web-chat/internal/hub"
	"web-chat/internal/logger"
	"web-chat/internal/middleware"
	"web-chat/internal/repository"
	"web-chat/internal/service"
	"web-chat/internal/worker"
)

func main() {
	cfg := logger.Config{
		Level:  logger.LevelFromString(os.Getenv("LOG_LEVEL")),
		Format: logger.FormatFromString(os.Getenv("LOG_FORMAT")),
	}

	log := logger.NewLogger(cfg)
	slog.SetDefault(log)

	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		host = "localhost"
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	db := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		host,
		os.Getenv("POSTGRES_DB"),
	)

	pool, err := repository.NewPool(context.Background(), db)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	onlineRedis, err := repository.NewRedisOnline(redisAddr)
	if err != nil {
		slog.Error("failed to connect to redis", "error", err)
		os.Exit(1)
	}

	defer onlineRedis.Close()

	userRepo := repository.NewUserPG(pool)
	roomRepo := repository.NewRoomPG(pool)
	msgRepo := repository.NewMessagePG(pool)

	userService := service.NewUserMemory(userRepo)
	roomService := service.NewRoomMemory(roomRepo)
	messageService := service.NewMessageMemory(msgRepo, roomRepo)

	userHandler := handler.NewUserHandler(userService)
	roomHandler := handler.NewRoomHandler(roomService)
	messageHandler := handler.NewMessageHandler(messageService)

	authHandler := handler.NewAuthHandler(userService)

	chatHub := hub.NewHub(onlineRedis)
	go chatHub.Run()

	chatPool := worker.NewPool(5)

	wsHandler := ws.ServeWS(chatHub, messageService, userService, chatPool)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("POST /auth/login", authHandler.Login)

	mux.Handle("GET /users/{id}", auth.AuthMiddleware(http.HandlerFunc(userHandler.GetByID)))
	mux.Handle("PUT /users/{id}", auth.AuthMiddleware(http.HandlerFunc(userHandler.Update)))
	mux.Handle("DELETE /users/{id}", auth.AuthMiddleware(http.HandlerFunc(userHandler.Delete)))

	mux.Handle("POST /rooms", auth.AuthMiddleware(http.HandlerFunc(roomHandler.Create)))
	mux.Handle("GET /rooms/{id}", auth.AuthMiddleware(http.HandlerFunc(roomHandler.GetByID)))
	mux.Handle("GET /rooms", auth.AuthMiddleware(http.HandlerFunc(roomHandler.GetAll)))
	mux.Handle("PUT /rooms/{id}", auth.AuthMiddleware(http.HandlerFunc(roomHandler.Update)))
	mux.Handle("DELETE /rooms/{id}", auth.AuthMiddleware(http.HandlerFunc(roomHandler.Delete)))

	mux.Handle("POST /rooms/{room_id}/messages", auth.AuthMiddleware(http.HandlerFunc(messageHandler.Create)))
	mux.Handle("GET /messages/{id}", auth.AuthMiddleware(http.HandlerFunc(messageHandler.GetByID)))
	mux.Handle("GET /rooms/{room_id}/messages", auth.AuthMiddleware(http.HandlerFunc(messageHandler.GetByRoomID)))
	mux.Handle("DELETE /messages/{id}", auth.AuthMiddleware(http.HandlerFunc(messageHandler.Delete)))
	mux.Handle("DELETE /rooms/{room_id}/messages", auth.AuthMiddleware(http.HandlerFunc(messageHandler.DeleteByRoom)))

	mux.Handle("GET /ws/chat/{room_id}", auth.AuthMiddleware(wsHandler))

	rateLimiter := middleware.NewRateLimiter(10)
	muxWithLogging := middleware.LoggingMiddleware(mux)
	muxLogRate := rateLimiter.Middleware(muxWithLogging)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: muxLogRate,
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	go http.ListenAndServe(":6060", nil)
	go func() {
		slog.Info("server starting", "port", 8080)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "error", err)
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down...")
	chatPool.Shutdown()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	srv.Shutdown(shutdownCtx)
	onlineRedis.Close()
	pool.Close()
	slog.Info("server stopped")
}
