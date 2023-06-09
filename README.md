## Project Setup

Before proceeding with the project setup, please ensure that you have Docker installed on your machine. If you haven't installed Docker yet, follow the instructions below:

Visit the official Docker website at https://www.docker.com/ and navigate to the Downloads section.
Choose the appropriate version of Docker for your operating system (Windows, macOS, or Linux).
Download and run the Docker installer.
Once Docker is successfully installed, you'll be able to proceed with the project setup and utilize its containerization capabilities.

```bash
git clone https://github.com/Wave-95/boards.git

cd boards
```

To start the project and bring up all the necessary containers, execute the following:

```bash
docker-compose up
```

This will set up the database (:5432), server (:8080), and frontend (:3000) containers. The database migrations will be run when the container is created for the server.

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