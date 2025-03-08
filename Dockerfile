FROM golang:1.21-alpine AS builder

ARG SERVICE=api

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -a -o /app/bin/${SERVICE} ./cmd/${SERVICE}

FROM alpine:latest

ARG SERVICE=api

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/bin/${SERVICE} /app/
COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /app/configs /app/configs

RUN adduser -D -g '' appuser \
    && chown -R appuser:appuser /app

USER appuser

CMD ["/app/${SERVICE}"] 