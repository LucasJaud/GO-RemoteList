# === etapa de build ===
FROM golang:1.21 AS builder
WORKDIR /app

# copia go.mod e go.sum primeiro (cache eficaz)
COPY go.mod ./
RUN go mod download

# copia o restante do código
COPY . .
RUN go build -o app

# === etapa de runtime enxuta ===
FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/app .
EXPOSE 8080
CMD ["./app"]
