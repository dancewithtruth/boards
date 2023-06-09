## Project Setup

This application is containerized with Docker. Please make sure you have docker installed. To set up the project, first clone the source code:

```bash
git clone https://github.com/Wave-95/boards.git

cd boards
```

Then run the `docker-compose.yml` file:

```bash
docker-compose up
```

This will set up the database (:5432), server (:8080), and frontend (:3000) containers. The database migrations will be run when the container is run for the server image.

You may need to run `chmod +x ./start.sh` locally since your host machine is mounted onto the container and will need access to the `start.sh` script to start the backend server.

## Development

The `docker-compose.yml` file loads in env vars for the server and frontend containers in their respective `.env` files. Depending on the `ENV` env varaible, the containers will either boot up in development or production mode. The server uses [air](https://github.com/cosmtrek/air) and the frontend uses `next dev`. 

## Database Migrations

Database migrations are run using [`golang-migrate`](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate).

### Running migrations

Access the server container first by finding the container ID

```bash
docker ps
```

```bash
docker exec -it <SERVER_CONTAINER_ID> /bin/bash
```

Once inside the container's terminal, run the following migrate command


```bash
migrate -path ./db/migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@db:${DB_PORT}/${DB_NAME}?sslmode=disable" up
```

To run the down migrations, simply change the `up` command to `down`.

### Creating migrations

Inside the docker container, run

```bash
migrate create -ext sql -dir db/migrations -seq <MIGRATION_NAME>
```