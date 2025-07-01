FROM golang:latest

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

COPY . .

RUN go build -o auth-service ./cmd/main.go

CMD ["./auth-service"]
