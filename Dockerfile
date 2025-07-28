# Билдер-этап
FROM golang:1.24.5-alpine AS builder
RUN apk add --no-cache git 
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /app/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /gin-app  

# Финальный образ
FROM alpine:latest
RUN apk add --no-cache tzdata ca-certificates  
WORKDIR /app
COPY --from=builder /gin-app .
COPY .env .  
EXPOSE 8081  

CMD ["./gin-app"]