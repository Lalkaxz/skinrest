GOOSE_DRIVER=postgres
GOOSE_DBSTRING="postgresql://postgres:0000@localhost:5432/skinRestDB?sslmode=disable"
build:
	docker compose build --no-cache

run-local:
	export SERVER_HOST="localhost" 
	export SERVER_PORT="8081" 
	export GIN_MODE="debug" 
	export API_ENV="local" 
	export DATABASE_DRIVER="postgres" 
	export DATABASE_HOST="localhost" 
	export DATABASE_PORT=5432 
	export DATABASE_USER="postgres" 
	export DATABASE_PASSWORD="0000" 
	export DATABASE_NAME="skinRestDB" 
	export DATABASE_SSL="disable" 
	export AUTH_JWT_SECRET="8ddeefb1f8c17f17864b0512c5148319848614a11efaed0b247c5cb2e19122e2" 
	go run ./cmd/server/main.go

up:
	docker compose up

down:
	docker compose down

restart:
	docker compose restart

migration-up:
	goose -dir migrations $(GOOSE_DRIVER) $(GOOSE_DBSTRING) up

migration-down:
	goose -dir migrations $(GOOSE_DRIVER) $(GOOSE_DBSTRING) down
