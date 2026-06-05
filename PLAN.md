# Chat Backend - Проектная Спецификация

## 📋 Database Schema

## 🔌 REST API Endpoints

### **User Management**

| Метод  | Endpoint             | Описание                           |
| ------ | -------------------- | ---------------------------------- |
| POST   | `/users`             | Создать пользователя               |
| GET    | `/users/{id}`        | Получить информацию о пользователе |
| PUT    | `/users/{id}`        | Обновить информацию пользователя   |
| DELETE | `/users/{id}`        | Удалить пользователя               |
| GET    | `/users/{id}/status` | Получить статус (онлайн/офлайн)    |

### **Room Management**

| Метод  | Endpoint                           | Описание                                   |
| ------ | ---------------------------------- | ------------------------------------------ |
| POST   | `/rooms`                           | Создать комнату                            |
| GET    | `/rooms`                           | Получить список всех комнат                |
| GET    | `/rooms/{id}`                      | Получить информацию о комнате              |
| PUT    | `/rooms/{id}`                      | Обновить комнату                           |
| DELETE | `/rooms/{id}`                      | Удалить комнату (только создатель)         |
| GET    | `/rooms/{id}/online-users`         | Получить онлайн пользователей в комнате    |
| GET    | `/rooms/{room_id}/users/{user_id}` | Получить информацию о юзере в этой комнате |

### **Messages**

| Метод  | Endpoint                                                  | Описание                                           |
| ------ | --------------------------------------------------------- | -------------------------------------------------- |
| GET    | `/rooms/{room_id}/messages?limit={limit}&offset={offset}` | История сообщений (paginated)                      |
| GET    | `/messages/{message_id}`                                  | Получить одно сообщение                            |
| DELETE | `/messages/{message_id}`                                  | Удалить сообщение (только автор)                   |
| DELETE | `/rooms/{id}/messages`                                    | Удалить все сообщения в комнате (только создатель) |

---

## 📡 WebSocket API

### **Connection**

```
ws://localhost:8080/ws/chat/{room_id}?user_id={user_id}
```

### **Message Types - From Client**

**Send Message:**

```json
{
    "type": "message",
    "text": "Hello everyone!",
    "room_id": "123"
}
```

### **Message Types - From Server (Broadcast)**

**New Message:**

```json
{
    "type": "message",
    "id": "456",
    "user_id": "789",
    "username": "john",
    "text": "Hello everyone!",
    "room_id": "123",
    "created_at": "2026-05-21T12:30:45Z"
}
```

**User Joined:**

```json
{
    "type": "user_joined",
    "user_id": "789",
    "username": "john",
    "room_id": "123",
    "online_count": 5
}
```

**User Left:**

```json
{
    "type": "user_left",
    "user_id": "789",
    "username": "john",
    "room_id": "123",
    "online_count": 4
}
```

**Error:**

```json
{
    "type": "error",
    "message": "Invalid message format"
}
```

---

## 🏗️ Как выстраивать папки

### **Основной принцип:**

Go проекты в production делятся на **слои архитектуры**:

```
User Request
    ↓
Handler (REST)
    ↓
Service (Business Logic)
    ↓
Repository (Database)
    ↓
Database
```

Каждый слой = своя папка в `internal/`

### **Что где лежит:**

**`cmd/`** - Entry point (откуда стартует приложение)

- Только `main.go`
- Инициализация всего
- Запуск сервера

**`internal/domain/`** - Models (сущности)

- User struct
- Room struct
- Message struct
- Никакой логики, только data structures

**`internal/handler/` или `internal/delivery/http/`** - HTTP handlers

- Обработка REST запросов
- Парсинг URL, query params
- Вызов service layer
- Возврат JSON ответов

**`internal/delivery/websocket/`** - WebSocket handler

- Upgrade connection
- Читать/писать WebSocket frames
- Вызов hub или service

**`internal/service/`** - Business Logic

- User creation logic
- Room creation logic
- Message sending logic
- Это где "мозг" приложения

**`internal/repository/`** - Database abstraction

- Query builders
- CRUD операции
- Connection pooling
- Ничего про бизнес-логику

**`internal/hub/`** - Message hub

- Broadcasting
- Connection management (add/remove)
- Goroutine management

**`internal/logger/`** - Logging

- slog setup
- Structured logging wrapper

**`internal/cache/` или `internal/redis/`** - Redis wrapper

- Connection to Redis
- Pub/Sub wrapper
- Caching operations

**`internal/middleware/`** - HTTP middleware

- Auth check
- Logging
- Error handling

**`pkg/`** - Reusable packages

- Если есть util функции что используются везде
- Обычно пусто на старте

### **Root level files:**

- `docker-compose.yaml` - DB + Redis
- `Dockerfile` - Application container
- `.dockerignore` - Исключить файлы из Docker build
- `.env` - Environment variables
- `Makefile` - Commands (make run, make test)
- `go.mod` / `go.sum` - Dependencies
- `README.md` - Documentation

---

## 🎯 Как создавать в правильном порядке

**Шаг 1:** Создай root папку

```bash
mkdir chat-backend
cd chat-backend
go mod init chat-backend
```

**Шаг 2:** Создай основные папки

```bash
mkdir cmd
mkdir internal
```

**Шаг 3:** Внутри `internal/` создай слои **в порядке dependency**

```bash
mkdir internal/domain        # Models (no dependencies)
mkdir internal/repository    # Database (depends on domain)
mkdir internal/service       # Business logic (depends on repository + domain)
mkdir internal/handler       # HTTP (depends on service)
mkdir internal/delivery
mkdir internal/delivery/websocket  # WebSocket (depends on service + hub)
mkdir internal/hub           # Broadcasting (depends on domain)
mkdir internal/logger        # Logging
mkdir internal/cache         # Redis
mkdir internal/middleware    # Middleware
```

**Шаг 4:** Создай files в `cmd/`

```bash
touch cmd/main.go
```

**Шаг 5:** Создай config files в root

```bash
touch docker-compose.yaml
touch Dockerfile
touch .dockerignore
touch .env
touch Makefile
touch README.md
```

---

## 💡 Важные правила

**1. Dependency Direction:**

- `cmd/` знает все
- `handler/` знает `service/`, но `service/` НЕ знает `handler/`
- `service/` знает `repository/`, но `repository/` НЕ знает `service/`
- `domain/` никого не знает (only data)

**2. Структура папок отражает зависимости:**

```
Если тебе нужна функция из handler в service
→ Это НЕПРАВИЛЬНО (нарушается архитектура)

Если тебе нужна функция из service в handler
→ Это ПРАВИЛЬНО (правильное направление)
```

**3. Каждая папка может иметь несколько файлов:**

```
internal/handler/
  ├── user_handler.go
  ├── room_handler.go
  └── message_handler.go

internal/service/
  ├── user_service.go
  ├── room_service.go
  └── message_service.go

internal/repository/
  ├── user_repository.go
  ├── room_repository.go
  └── message_repository.go
```

**4. Interfaces в каждом слое:**

```go
// repository/user_repository.go
type UserRepository interface {
  GetByID(ctx, id) (*User, error)
  Create(ctx, user) error
}

// service/user_service.go
type UserService struct {
  repo UserRepository  // Inject interface, not concrete type
}
```

---

## 🎬 После создания структуры

Когда создашь все папки:

1. Создай empty `go.mod`
2. Создай `docker-compose.yaml` (PostgreSQL + Redis)
3. Создай `Dockerfile`
4. Создай `Makefile`

**Потом пиши мне структуру** → я проверю правильно ли выстроил!

---

## 📝 Implementation Phases

### **Phase 0: Setup (День 1)**

- [ ] `go mod init chat-backend`
- [ ] Создать структуру папок
- [ ] Написать docker-compose.yaml (PostgreSQL + Redis)
- [ ] Написать Dockerfile
- [ ] Написать .env файл
- [ ] Написать Makefile

### **Phase 1: HTTP Server + WebSocket (Дни 2-4)**

- [ ] HTTP мux и роутинг
- [ ] WebSocket upgrade
- [ ] Подключение/отключение
- [ ] Отправка и получение сообщений
- [ ] Context и cancellation
- [ ] Логирование

### **Phase 2: Hub & Broadcasting (Дни 5-6)**

- [ ] Hub структура
- [ ] Регистрация/удаление connections
- [ ] Message broadcast (fan-out)
- [ ] sync.Mutex protection
- [ ] Graceful shutdown

### **Phase 3: PostgreSQL (Дни 7-8)**

- [ ] Создать schema
- [ ] Connection pooling
- [ ] Repository layer
- [ ] CRUD операции
- [ ] Сохранение сообщений

### **Phase 4: Redis (Дни 9-10)**

- [ ] Redis connection
- [ ] Pub/Sub broadcast
- [ ] User presence tracking
- [ ] Message caching

### **Phase 5: Performance (Дни 11-12)**

- [ ] Load test (10k+ connections)
- [ ] pprof profiling
- [ ] Memory optimization
- [ ] Benchmark

### **Phase 6: Testing (Дни 13-14)**

- [ ] Unit tests
- [ ] Integration tests
- [ ] Error handling
- [ ] Documentation

---

## 🛠️ Technologies

| Категория | Инструмент             |
| --------- | ---------------------- |
| Language  | Go 1.23+               |
| Database  | PostgreSQL             |
| Cache     | Redis                  |
| WebSocket | gorilla/websocket      |
| HTTP      | net/http               |
| UUID      | github.com/google/uuid |
| Testing   | testing (stdlib)       |
| Profiling | runtime/pprof          |

---

## 🚀 Что делать теперь

1. Создай модуль: `go mod init chat-backend`
2. Создай структуру папок (как выше)
3. Напиши `docker-compose.yaml` (PostgreSQL + Redis)
4. Напиши `Dockerfile`
5. Напиши `.env` файл
6. Напиши `Makefile` с командами:
    - `make run` - запустить приложение
    - `make test` - запустить тесты
    - `make docker-up` - docker-compose up
    - `make docker-down` - docker-compose down

**Потом пиши мне что создал** - проверю правильно ли или переделать надо!

---

## 💡 Важные моменты

- **Все queries должны быть prepared statements** (security + performance)
- **Все операции с goroutines должны использовать context** (cancellation)
- **Все логирование структурированное JSON** (slog)
- **Нет глобальных переменных** (dependency injection)
- **Interface-based design** (easy testing)
