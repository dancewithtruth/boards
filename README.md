## Boards
Boards is a collaborative desktop web application for creating and organizing notes on a digital canvas. Boards is a demo project showcasing event-driven microservices written in Golang as well as a stateful, interactive UI built with Next.js, Tailwind, and [React DnD](https://react-dnd.github.io/react-dnd/about). The containerized applicaton is orchestrated using Kubernetes and hosted on AWS. Perhaps the most interesting aspect of the application is how data is synchronized in real-time across client connections using Redis Pub-Sub and WebSocket events (create, move, edit posts).

<img src="frontend/public/Hero.png" alt="Boards" width="500"/>

## Technology Stack
#### Frontend
- [Next.js 13](https://nextjs.org/)
- [Tailwind CSS](https://tailwindcss.com/)
- [React Drag and Drop](https://react-dnd.github.io/react-dnd/about)
#### Backend
- [Chi Router](https://github.com/go-chi/chi)
- [Gorilla WebSockets](https://github.com/gorilla/websocket)
- [RabbitMQ](https://www.rabbitmq.com/tutorials/tutorial-one-go.html)
- [Redis PubSub](https://redis.io/docs/interact/pubsub/)
#### Storage / DB
- PostgreSQL
- [sqlc](https://sqlc.dev/)
- [pgx](https://github.com/jackc/pgx)
- [Golang Migrate v4](https://github.com/golang-migrate/migrate)
#### Infrastructure / Cloud
- Docker
- Kubernetes
- AWS: EKS, EC2, Route 53

## Architecture Overview
![architecture](docs/architecture.svg)

## Services

No. | Service | Local | Hosted
--- | --- | --- | ---
1 | backend-core | [http://localhost:8080](http://localhost:8080) | [https://api.useboards.com](https://api.useboards.com)
2 | backend-notification | [http://localhost:8082](http://localhost:8082) |
3 | web | [http://localhost:3000](http://localhost:3000) | [https://useboards.com](https://useboards.com)
4 | docs | [http://localhost:8081](http://localhost:8081) | [https://docs.useboards.com](https://docs.useboards.com)

## Relevant Blog Posts
- [WebSockets Client and Hub Architecture](https://medium.com/@wu.victor.95/building-a-go-websocket-for-a-live-collaboration-tool-pt-1-f7e5374b1f47)
- [Implementing Go WebSockets using TDD](https://medium.com/@wu.victor.95/building-a-go-websocket-for-a-live-collaboration-tool-pt-2-5728cd6ec801)
- [Implementing User Authentication Event](https://medium.com/@wu.victor.95/building-a-go-websocket-for-a-live-collaboration-tool-pt-3-b9a6b23f7fef)
- [Hashing Passwords and Authenticating Users](https://medium.com/@wu.victor.95/hashing-passwords-and-authenticating-users-with-bcrypt-dc2fdd978568)
- [Building the Posts API](https://medium.com/@wu.victor.95/building-a-post-service-for-our-websocket-endpoint-using-clean-architecture-tdd-f39aae9b2041)
- [Building Frontend with Next.js & Tailwind](https://medium.com/@wu.victor.95/intro-8435223725f0)
- [Fuzzy Email Search](https://medium.com/@wu.victor.95/new-feature-invite-members-to-a-board-cddfb6657131)
- [E2E Implementation for User Invites](https://medium.com/@wu.victor.95/new-feature-board-invitations-pt-2-549e071d0338)
- [Setting up Kubernetes](https://medium.com/@wu.victor.95/deploying-with-kubernetes-d3a9e9aad767)
- [Deploying Application Cluster to AWS](https://medium.com/@wu.victor.95/deploying-application-to-aws-1d9b4e758de0)
- [Frontend: Grouping Posts](https://medium.com/@wu.victor.95/boards-new-feature-grouping-posts-pt-1-680a98701c9b)
- [Frontend: Ordering Posts](https://medium.com/@wu.victor.95/boards-new-feature-ordering-posts-e3984adcdef5)
- [Making Backend Server Stateless](https://medium.com/@wu.victor.95/stateless-websocket-server-using-redis-pubsub-bf5f70435ba0)
- [Creating Email Notification Service using RabbitMQ](https://medium.com/@wu.victor.95/creating-a-notification-service-using-rabbitmq-a488c3d5b8bf)
- [Implementing Email Verifications](https://medium.com/@wu.victor.95/sending-email-verifications-with-smtp-pt-1-c4ededf5442a)
- [Using SMTP to Send Emails](https://medium.com/@wu.victor.95/handling-rabbitmq-tasks-and-sending-emails-with-smtp-d2fd6bac695e)

## Project Setup

Before proceeding with the project setup, please ensure that you have [Docker](https://www.docker.com/) installed on your machine. 

```bash
git clone https://github.com/Wave-95/boards.git

cd boards

docker-compose up
```

This will start up the frontend, backend, and notification service as well as the PostgreSQL, Redis, and RabbitMQ servers on your local machine. Database migrations should automatically run when a container is created for the backend service during docker compose. If you would like to insert test data into the database, use `make testdata`. 

## Development

The `docker-compose.yml` file loads in env vars for the server and frontend containers in their respective `.env` files. Depending on the `ENV` env varaible, the containers will either boot up in development or production mode. The backend uses [air](https://github.com/cosmtrek/air) and the frontend uses `next dev`. Feel free to rename the `.env.example` files to `.env` to get local development up and running.

## Database Migrations

Database migrations are run using [`golang-migrate`](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate). The `golang-migrate` dependency should already be available in the backend container. You can either run migrations manually via an interactive shell or use the `make` commands:

```bash
make migrate-up

make migrate-down

make testdata

make migrate-create
```

## Run Tests

```bash
make test
```

## To Do
- WebSocket API documentation
- Implement mobile responsiveness
- Implement high throughput data storage solution
- Refactor frontend components
- Add post voting
- Add board duplication
