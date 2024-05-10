FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o app ./cmd/user-api-service/main.go

EXPOSE 8080

ENV CONFIG_PATH="./config/local.yml"

CMD ["./app"]
