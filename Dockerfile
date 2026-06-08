# Estágio 1: Compilação do código Go
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copia os arquivos de dependências
COPY go.mod go.sum ./
RUN go mod download

# Copia o código-fonte restante
COPY main.go ./

# Compila o binário estaticamente para o Linux
RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go

# Estágio 2: Imagem final leve para execução
FROM alpine:latest
WORKDIR /root/

# Copia o binário gerado no estágio de compilação
COPY --from=builder /app/main .

EXPOSE 8080
CMD ["./main"]