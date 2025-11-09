# Set up VSCode environment
init-vscode:
	cp .env.sample .env
	cd .vscode && cp settings.json.sample settings.json
# Set up database
set-up-db:
	docker compose down && \
	docker compose up nexinvent-db -d --wait
	go run migrations/clear_database.go
# Runners
run:
	cd src/server && go run main.go