FROM golang:1.23

WORKDIR /app

# Копируем файлы go.mod и go.sum
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем весь код из директории app и других папок
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/docker-gs-ping ./app/main.go

# Открываем порт 8000
EXPOSE 8000

# Запускаем приложение
CMD ["/app/docker-gs-ping"]