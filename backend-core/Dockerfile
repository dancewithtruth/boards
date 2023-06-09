FROM golang:latest

# Install air for development hot reloading
RUN go install github.com/cosmtrek/air@latest

# Intall go migrate tool
# https://github.com/golang-migrate/migrate/tree/master/cmd/migrate
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

WORKDIR /app

COPY go.* ./
RUN go mod download && go mod verify

COPY . .

RUN go build -o ./cmd/main ./cmd/main.go

EXPOSE 8080

CMD [ "./start.sh" ] 