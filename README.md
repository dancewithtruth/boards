## Project Setup

`docker-compose up`

You may need to run `chmod +x ./start.sh` locally since your host machine is mounted onto the container and will need access to the `start.sh` script to start the backend server.

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