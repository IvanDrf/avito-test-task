FROM golang:1.25.1-alpine

RUN apk add --no-cache git gcc musl-dev sqlite

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -o app ./cmd/app/main.go
RUN CGO_ENABLED=1 go build -o migrator ./cmd/migrator/main.go

RUN mkdir -p storage/sqlite

EXPOSE 8080

CMD ["sh", "-c", "./migrator --storage-path=./storage/sqlite/pr.db --migrations-path=./internal/migrations && ./app"]