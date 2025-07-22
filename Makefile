up:
	@docker compose up -d --build

down:
	docker compose down

logs-app-1:
	docker-compose logs -f rinha-api-1