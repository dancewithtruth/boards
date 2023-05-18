## Project Setup

Make sure you have docker installed:

`docker-compose up`

You may need to run `chmod +x ./start.sh` locally since your host machine is mounted onto the container and will need access to the `start.sh`

## Migrations

Identify the container ID or name (should be `boards-server`)

`docker ps`

Access the container's terminal

`docker exec -it <container_id_or_name> /bin/bash`

From here you can run commands inside the container. The container already has `golang-migrate` installed. To add a new migration,

`migrate create -ext sql -dir db/migrations -seq create_users_table`