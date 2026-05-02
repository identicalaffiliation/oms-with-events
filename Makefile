.PHONY: up down

default: up

up:
	docker compose up -d --build

down:
	docker compose down -v