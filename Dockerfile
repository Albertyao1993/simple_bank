# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go 
RUN apk add --no-cache curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
RUN mv migrate /usr/bin/migrate

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /usr/bin/migrate /usr/bin/migrate
COPY --from=builder /app/db/migration ./db/migration
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
RUN chmod +x start.sh

CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]
