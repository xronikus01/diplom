# Advanced Blog Management System

Дипломный проект на Go: REST API для блог-платформы с регистрацией пользователей, JWT-аутентификацией, управлением постами и комментариями, PostgreSQL и Docker.

## Функциональность

### Пользователи
- регистрация нового пользователя;
- вход пользователя;
- хеширование пароля через bcrypt;
- JWT-аутентификация.

### Посты
- создание поста;
- получение списка опубликованных постов;
- получение поста по ID;
- получение постов конкретного пользователя;
- обновление поста;
- удаление поста;
- отложенная публикация постов.

### Комментарии
- создание комментария;
- получение комментариев поста;
- обновление комментария;
- удаление комментария.

### Сервис
- health-check endpoint;
- middleware аутентификации;
- middleware логирования;
- worker для публикации запланированных постов;
- graceful shutdown.

---

## Стек

- Go 1.22
- Chi Router
- PostgreSQL 15
- JWT (`github.com/golang-jwt/jwt/v5`)
- bcrypt (`golang.org/x/crypto/bcrypt`)
- godotenv
- Docker / Docker Compose

---

## Структура проекта

```text
blog-api/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── handler/
│   │   ├── auth_handler.go
│   │   ├── comment_handler.go
│   │   ├── health.go
│   │   ├── helpers.go
│   │   └── post_handler.go
│   ├── middleware/
│   │   ├── auth.go
│   │   └── logging.go
│   ├── model/
│   │   └── models.go
│   ├── repository/
│   │   ├── comment_repo.go
│   │   ├── interfaces.go
│   │   ├── post_repo.go
│   │   └── user_repo.go
│   ├── service/
│   │   ├── comment_service.go
│   │   ├── errors.go
│   │   ├── post_service.go
│   │   └── user_service.go
│   └── worker/
│       └── scheduler.go
├── migrations/
│   ├── 001_init_schema.sql
│   └── 002_add_indexes.sql
├── pkg/
│   ├── auth/
│   │   ├── jwt.go
│   │   └── password.go
│   ├── config/
│   │   └── config.go
│   └── database/
│       └── db.go
├── .env.example
├── .gitignore
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
```

## Переменные окружения

Файл `.env`:

```env
# Database
DB_HOST=127.0.0.1
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=blog_db
DB_SSLMODE=disable

# JWT
JWT_SECRET=super-secret-key

# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# Environment
ENV=development
```

---

## Запуск проекта

### 1. Установить зависимости

```bash
go mod tidy
```

### 2. Запустить PostgreSQL

```bash
docker compose up -d db
```

### 3. Применить миграции

Подключиться к БД:

```bash
docker compose exec db psql -U postgres -d blog_db
```

Создать таблицы:

```sql
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    author_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'published',
    publish_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY,
    post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    author_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

Создать индексы:

```sql
CREATE INDEX IF NOT EXISTS idx_posts_author_id ON posts(author_id);
CREATE INDEX IF NOT EXISTS idx_posts_status ON posts(status);
CREATE INDEX IF NOT EXISTS idx_posts_publish_at ON posts(publish_at);
CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments(post_id);
CREATE INDEX IF NOT EXISTS idx_comments_author_id ON comments(author_id);
```

Выйти из `psql`:

```sql
\q
```

### 4. Запустить приложение

```bash
go run ./cmd/api/main.go
```

Сервер будет доступен по адресу:

```text
http://localhost:8080
```

---

## Запуск через Docker

### Запуск PostgreSQL

```bash
docker compose up -d db
```

### Сборка и запуск приложения

```bash
docker build -t blog-api .
docker run --rm -p 8080:8080 --env-file .env blog-api
```

---

## API эндпоинты

### Публичные

#### Health check
- `GET /api/health`

#### Пользователи
- `POST /api/register`
- `POST /api/login`

#### Посты
- `GET /api/posts`
- `GET /api/posts/{id}`
- `GET /api/users/{id}/posts`

#### Комментарии
- `GET /api/posts/{postId}/comments`

### Защищённые

Требуют заголовок:

```text
Authorization: Bearer <token>
```

#### Посты
- `POST /api/posts`
- `PUT /api/posts/{id}`
- `DELETE /api/posts/{id}`

#### Комментарии
- `POST /api/posts/{postId}/comments`
- `PUT /api/comments/{id}`
- `DELETE /api/comments/{id}`

---

## Примеры запросов

### Health

```bash
curl http://localhost:8080/api/health
```

### Регистрация

```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "username":"stepan",
    "email":"stepan@example.com",
    "password":"Password123"
  }'
```

### Логин

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email":"stepan@example.com",
    "password":"Password123"
  }'
```

### Создание поста

```bash
curl -X POST http://localhost:8080/api/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TOKEN" \
  -d '{
    "title":"Мой первый пост",
    "content":"Это содержимое моего первого поста."
  }'
```

### Получение всех постов

```bash
curl http://localhost:8080/api/posts
```

### Получение поста по ID

```bash
curl http://localhost:8080/api/posts/1
```

### Создание комментария

```bash
curl -X POST http://localhost:8080/api/posts/1/comments \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TOKEN" \
  -d '{
    "content":"Отличный пост!"
  }'
```

### Получение комментариев поста

```bash
curl http://localhost:8080/api/posts/1/comments
```

---

## Что реализовано

- регистрация пользователей;
- логин пользователей;
- bcrypt-хеширование паролей;
- JWT-аутентификация;
- middleware авторизации;
- middleware логирования;
- CRUD для постов;
- CRUD для комментариев;
- PostgreSQL в Docker;
- worker для отложенной публикации постов;
- корректная обработка ошибок;
- JSON-ответы;
- graceful shutdown.

---

## Что проверено вручную

Проверены следующие сценарии:
- `GET /api/health`
- `POST /api/register`
- `POST /api/login`
- `POST /api/posts`
- `GET /api/posts`
- `GET /api/posts/{id}`
- `GET /api/users/{id}/posts`
- `PUT /api/posts/{id}`
- `DELETE /api/posts/{id}`
- `POST /api/posts/{postId}/comments`
- `GET /api/posts/{postId}/comments`
- `PUT /api/comments/{id}`
- `DELETE /api/comments/{id}`

---

## Комментарий по запуску на Windows

На Windows порт `5432` может быть занят локальным PostgreSQL, поэтому в проекте используется:

```env
DB_PORT=5433
```

и в `docker-compose.yml`:

```yaml
ports:
  - "5433:5432"
```

---

## Автор

Гаран Степан
