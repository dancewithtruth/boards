FROM golang:latest

# Install air for development hot reloading
RUN go install github.com/cosmtrek/air@latest

# Install dlv for debugging
RUN go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /app/src

COPY go.* ./
RUN go mod download && go mod verify

COPY . .

RUN go build -o /app/main ./cmd/main.go

EXPOSE 2346

CMD [ "./start.sh" ] 