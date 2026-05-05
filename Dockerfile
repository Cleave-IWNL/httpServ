FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -tags musl -o server ./cmd/server

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/server .
COPY migrations/ migrations/

EXPOSE 8080

CMD ["./server"]