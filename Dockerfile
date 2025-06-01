FROM golang:1.24.3-alpine

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0

RUN go build -o property-list ./cmd

COPY .env .

EXPOSE 8080

CMD ["./property-list"]
