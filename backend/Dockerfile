# Etapa de compilação
FROM golang:1.25-alpine AS builder
# Defina o diretório de trabalho dentro do contêiner
WORKDIR /app
# Copie o módulo go e baixe as dependências
COPY backend/go.mod ./
COPY backend/go.sum ./
RUN go mod download
RUN apk --update add ca-certificates
# Copie o código fonte para o diretório de trabalho
COPY backend/ .
# Compile o aplicativo Go para um binário estático.
# Note que CGO está desabilitado para garantir um binário estático
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd

# Etapa de execução
FROM scratch
# Copie o binário compilado da etapa de compilação
COPY --from=builder /app/main .
# Copie o arquivo de certificado SSL para o contêiner
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
# Copia templates para o contêiner
COPY backend/templates /templates
# Comando para executar o binário
CMD ["./main"]
