## Project Setup

Please make sure you have docker installed:

`docker-compose up`

You may need to run `chmod +x ./start.sh` locally since your host machine is mounted onto the container and will need access to the `start.sh` script to start the backend server.

## Migrations

First access the container's terminal

`docker exec -it boards-server /bin/bash`

Now that you have access to the installed `golang-migrate` library, you can run migrations in the container:

`migrate create -ext sql -dir db/migrations -seq create_users_table`