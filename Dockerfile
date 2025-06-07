# Use a imagem base oficial do Go
FROM golang:1.21-alpine

# Defina o diretório de trabalho
WORKDIR /app

# Copie os arquivos go.mod e go.sum (se existirem)
COPY go.mod go.sum ./

# Baixe as dependências
RUN go mod download

# Copie o código fonte
COPY . .

# Compile a aplicação
RUN go build -o main .

# Exponha a porta (se for uma aplicação web)
EXPOSE 8080

# Execute a aplicação
CMD ["./main"]