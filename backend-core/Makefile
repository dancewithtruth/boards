include .env
CONTAINER_NAME_BACKEND := backend-core
CONTAINER_NAME_DB := db-core
DSN := postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable
MIGRATE := docker exec -it ${CONTAINER_NAME_BACKEND} migrate -path ./db/migrations -database "${DSN}"

.PHONY: migrate-up
migrate-up:
	@echo "Running all up migrations..."
	@${MIGRATE} up

.PHONY: migrate-down
migrate-down:
	@echo "Running all down migrations..."
	@${MIGRATE} down

.PHONY: testdata
testdata:
	@echo "Resetting database before insert..."
	@${MIGRATE} drop
	@${MIGRATE} up
	@echo "Inserting testdata..."
	@for file in testdata/*.sql; do \
        echo "Executing SQL file: $$file"; \
        docker exec -it ${CONTAINER_NAME_DB} psql "${DSN}" -f $$file; \
    done
	
.PHONY: migrate-create
migrate-create:
	@read -p "Enter a name for your migration file: " name; \
	${MIGRATE} create -ext sql -dir ./db/migrations $${name// /_}
	