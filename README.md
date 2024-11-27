# Medods


- [О проекте](#о-проетке)
- [Настройка](#настройка)
- [Запуск](#запуск)
- [Использование](#использование)


## О проетке

Сиситема авторизации пользователей через JWT токены с проверкой ip адреса

## Настройка

Для конфигурации проекта необходимо внести изменения в yaml файлы которые находятся в директории `/config`

Для замены файла конфигурации в проекте необходимо изменить путь к файлу в файле `cmd/main.go`

```go
cfg, err := config.MustLoad("config/docker.yaml")
```
## Запуск

### Локальный запуск

Для локального запуска проекте необходимо в корне проекта ввести команду

```bash
go run cmd/main.go
```
### Тесты

Для запуска тестов необходимо внести изменения в файл конфигурации, который располагается по пути `config/test.yaml`

После в корне проекта выполнить команду
```bash
go test -v ./tests/...
```

### Запуск через Docker

Для запуска через Docker необходимо в корне проекта выполнить команду

```bash
docker-compose up -d --build
```
Для остановки контейнера необходимо выполнить команду
```bash
docker-compose down -v
```

## Использование

Существует 3 HTTP энпоинта

### Создание пользователя

`POST /create`

Тело зопроса в формате JSON
```json
{
    "email": "testemail@example.com",
}
```

В ответе приходит JSON с данными о новом пользователе или ошибка.
Пример успешного ответа:
```json
{
    "detail": "new user was successfully created",
    "id": "986ae9e5-d02b-45cb-8328-94e0ea285e05"
}
```

### Создание токенов

`GET /tokens?userID=<ID пользователя>`

Пример запроса на получение пары токенов:

`http://localhost:8080/tokens?userID=986ae9e5-d02b-45cb-8328-94e0ea285e05`

В ответе приходит JSON с данными о новом пользователе или ошибка.
Пример успешного ответа:

```json
{
    "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzI3NDg4MTYsInVzZXJfaWQiOiI5ODZhZTllNS1kMDJiLTQ1Y2ItODMyOC05NGUwZWEyODVlMDUifQ.hC_YVYHQsYiNmJ-AB_hnrtrVkjcjl142d93yXMQHFI0g2U874ykgh3TXthuvmggUXJ8H-UptfXXpxq45QIxHqQ",
    "refresh_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzMzNTAwMTYsImlwX2FkZHJlc3MiOiIxOTIuMTY4LjY1LjE6NTM2NDUiLCJ1c2VyX2lkIjoiOTg2YWU5ZTUtZDAyYi00NWNiLTgzMjgtOTRlMGVhMjg1ZTA1In0.Umy683OQfHHEhdNsMR-xAq4CFMUQj7pYP6eE1LlAfHtr6SQDXgFxuQiOxeaVjGdNWGioawp4nJZXyPJFh2KK1A"
}
```

### Обновление токенов

`GET /refresh-tokens`

В заголовках запроса необходимо отправить refresh токен

`Authorization: Bearer <refresh token>`

В случае если запрос прошел успешно, в ответе вернется новая пара токенов:
```json
{
    "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzI3NDg4MTYsInVzZXJfaWQiOiI5ODZhZTllNS1kMDJiLTQ1Y2ItODMyOC05NGUwZWEyODVlMDUifQ.hC_YVYHQsYiNmJ-AB_hnrtrVkjcjl142d93yXMQHFI0g2U874ykgh3TXthuvmggUXJ8H-UptfXXpxq45QIxHqQ",
    "refresh_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzMzNTAwMTYsImlwX2FkZHJlc3MiOiIxOTIuMTY4LjY1LjE6NTM2NDUiLCJ1c2VyX2lkIjoiOTg2YWU5ZTUtZDAyYi00NWNiLTgzMjgtOTRlMGVhMjg1ZTA1In0.Umy683OQfHHEhdNsMR-xAq4CFMUQj7pYP6eE1LlAfHtr6SQDXgFxuQiOxeaVjGdNWGioawp4nJZXyPJFh2KK1A"
}
```