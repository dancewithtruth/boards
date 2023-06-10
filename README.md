## Project Setup

Before proceeding with the project setup, please ensure that you have [Docker](https://www.docker.com/) installed on your machine. 

```bash
git clone https://github.com/Wave-95/boards.git

cd boards

docker-compose up
```

This will start up the frontend (port 3000), backend (8080), and postgres database (5432) on your local machine. Database migrations should automatically run when the backend service is containerized during the compose step. If you would like test data to start, use `make testdata`. 

## Development

The `docker-compose.yml` file loads in env vars for the server and frontend containers in their respective `.env` files. Depending on the `ENV` env varaible, the containers will either boot up in development or production mode. The backend uses [air](https://github.com/cosmtrek/air) and the frontend uses `next dev`. 

## Database Migrations

Database migrations are run using [`golang-migrate`](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate). The `golang-migrate` dependency should already be available in the backend container. You can either run migrations manually via an interactive shell or use the `make` commands:

```bash
make migrate-up

make migrate-down

make testdata

make migrate-create
```