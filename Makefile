include .env

export PROJECT_ROOT=$(shell pwd)

webchat-run:
	@go run cmd/main.go