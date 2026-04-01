# Gotochat

Веб-чат на Go с регистрацией, JWT (Ed25519), WebSocket-комнатами и статической клиентской страницей. Сообщения хранятся в памяти комнат; учётные записи — в PostgreSQL. 

## Возможности

- Регистрация и вход по JSON API; после входа выдаются cookie `token` и `name`.
- Защищённый WebSocket: `GET /chat/{roomID}` (требуется валидный JWT в cookie `token`).
- Десять предсозданных комнат: `room0` … `room9`.
- Метрики Prometheus на `/metrics`; в `docker-compose` подключены Prometheus и Grafana.

## Переменные окружения

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `DATABASE_URL` | DSN PostgreSQL | `postgres://user:password@localhost:5432/gotodb` |

Сервер слушает `:8080` (задаётся в коде `cmd/main.go`).

## Запуск локально

1. Поднять PostgreSQL и создать БД (или использовать значения по умолчанию из таблицы выше).
2. Создать ключи в `./keys`.
3. При необходимости задать `DATABASE_URL`.
4. Из корня репозитория:

```bash
go run ./cmd/main.go
```

Открыть в браузере: `http://localhost:8080` — раздаётся `static/index.html`.

## Docker Compose

```bash
docker compose up --build
```

- Приложение: `http://localhost:8080`
- PostgreSQL: `localhost:5432` (user/password, БД `gotodb`)
- Prometheus: `http://localhost:9090`
- Grafana: `http://localhost:3000` (логин по умолчанию `admin`, пароль из `GF_SECURITY_ADMIN_PASSWORD` в `docker-compose.yaml`)

Убедитесь, что в образ попали файлы `keys/` при сборке, если нужны регистрация и вход.

## HTTP API

### `POST /register`

Тело (JSON):

```json
{ "username": "alice", "password": "secret" }
```

Успех: `201 Created`. Ошибка валидации или дубликат пользователя: `400`.

### `POST /login`

Тело (JSON):

```json
{ "username": "alice", "password": "secret" }
```

Успех: `200 OK`, устанавливаются cookie `token` и `name`. Неверные учётные данные: `401`.

### `GET /chat/{roomID}`

Требуется cookie `token` с валидным JWT. Обновление до WebSocket; для отображения имени в чате используется cookie `name`.

Комната должна существовать (например `room0`). Иначе: `404`.

### `GET /metrics`

Метрики в формате Prometheus (без JWT).

## Структура репозитория (кратко)

| Путь | Назначение |
|------|------------|
| `cmd/main.go` | Точка входа, подключение к БД, запуск HTTP |
| `internal/app` | Маршруты, middleware, сервисы |
| `internal/handler` | Обработчики register, login, chat |
| `internal/services` | Авторизация и логика чата |
| `internal/storage/postgres` | Пользователи, миграция таблицы `users` |
| `static/` | Статика (UI) |
| `deploy/` | Конфиги Prometheus и Grafana |
