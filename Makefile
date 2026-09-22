.PHONY: build up down logs test backend frontend

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f

backend-dev:
	cd backend && go run *.go

frontend-dev:
	cd frontend && npm run dev
