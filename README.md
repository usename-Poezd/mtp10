# Лабораторная работа №10 — Веб-разработка: FastAPI (Python) vs Gin (Go)

## Студент

- **ФИО**: Поярков Артём Алексеевич
- **Группа**: 221131

## Вариант 2

### Задания средней сложности

| # | Задание | Реализация |
|---|---------|-----------|
| М2 | Добавить middleware для логирования в Go | `go-service/internal/middleware/logger.go` — структурированный key=value лог |
| М4 | Создать FastAPI-сервис, который вызывает Go-сервис через HTTP | `python-service/app/` — прокси с Pydantic-валидацией и `/stats` |
| М6 | Сравнить скорость ответа FastAPI и Gin под нагрузкой (wrk) | `bench/run_benchmark.sh`, `BENCHMARK.md` |

### Задания повышенной сложности

| # | Задание | Реализация |
|---|---------|-----------|
| В2 | Создать API-шлюз на Go, маршрутизирующий запросы к микросервисам | `go-service/internal/gateway/` — `/api/go/*` и `/api/python/*` |
| В4 | WebSocket: чат на Go + подключение из Python | `go-service/internal/ws/` (Hub) + `python-service/ws_client.py` |

---

## Описание

Монорепозиторий с реализацией 5 заданий лабораторной работы №10 по сравнению подходов к веб-разработке на Python (FastAPI) и Go (Gin).

### Реализованные компоненты

| Сервис | Описание | Порт |
|--------|----------|------|
| **Go TODO API** | CRUD-сервис задач на Gin с in-memory хранилищем и структурированным logging middleware | 8080 |
| **FastAPI Proxy** | Python-сервис на FastAPI, проксирующий запросы к Go-сервису и добавляющий Pydantic-валидацию, поля `created_at`/`priority` и эндпоинт `/stats` | 8000 |
| **API Gateway** | Обратный прокси на Go/Gin, маршрутизирующий `/api/go/*` → Go-сервис и `/api/python/*` → FastAPI | 9000 |
| **WebSocket Chat** | Go-сервер с Hub-паттерном (gorilla/websocket) + Python CLI-клиент | 8081 |

---

## Технологии

- **Go** 1.22 — Gin, gorilla/websocket, net/http/httputil
- **Python** 3.12 — FastAPI, uvicorn, httpx, pydantic, websockets
- **Тестирование** — go test (с race detector), pytest, respx
- **Контейнеризация** — Docker, docker-compose
- **Нагрузочное тестирование** — wrk

---

## Структура проекта

```
lr10/
├── go-service/
│   ├── cmd/
│   │   ├── server/main.go       # Entry point TODO API (:8080)
│   │   ├── gateway/main.go      # Entry point API Gateway (:9000)
│   │   └── ws-server/main.go    # Entry point WebSocket (:8081)
│   ├── internal/
│   │   ├── app/app.go           # App bootstrap
│   │   ├── handlers/http/       # HTTP handlers (handler.go + todos.go)
│   │   ├── gateway/             # Gateway reverse proxy
│   │   ├── middleware/          # Structured logging middleware
│   │   ├── models/              # Todo model
│   │   ├── store/               # Thread-safe in-memory store
│   │   └── ws/                  # WebSocket Hub + handler
│   ├── Dockerfile
│   ├── Dockerfile.gateway
│   └── go.mod
├── python-service/
│   ├── app/
│   │   ├── main.py              # FastAPI endpoints
│   │   ├── models.py            # Pydantic models
│   │   └── service.py           # Go service HTTP client
│   ├── tests/
│   │   ├── test_todos.py
│   │   ├── test_stats.py
│   │   └── test_ws_client.py
│   ├── ws_client.py             # WebSocket CLI client
│   ├── Dockerfile
│   └── requirements.txt
├── bench/
│   └── run_benchmark.sh         # wrk benchmark script
├── docker-compose.yml
├── BENCHMARK.md                 # Load test results
└── README.md
```

---

## Сборка проекта

### Локальная сборка

**Go-сервисы:**
```bash
cd go-service
go build ./...          # Собрать все бинарники
go vet ./...            # Статический анализ (0 issues)
go test -race -v ./...  # Юнит-тесты с race detector
```

**Python-сервис:**
```bash
cd python-service
python3 -m venv venv
. venv/bin/activate
python3 -m pip install -r requirements.txt
ruff check app/ tests/  # Линтер (0 issues)
python3 -m pytest -v    # Юнит-тесты
```

### Docker сборка

```bash
docker compose build
```

---

## Запуск

### Запуск через Docker Compose (рекомендуется)

```bash
docker compose up -d
```

Сервисы запустятся в правильном порядке (go-service → python-service → gateway).

### Локальный запуск

**Go TODO API** (порт 8080):
```bash
cd go-service && go run cmd/server/main.go
```

**Python FastAPI** (порт 8000):
```bash
cd python-service && . venv/bin/activate && uvicorn app.main:app --port 8000
```

**API Gateway** (порт 9000):
```bash
cd go-service && go run cmd/gateway/main.go
```

**WebSocket сервер** (порт 8081):
```bash
cd go-service && go run cmd/ws-server/main.go
```

---

## Примеры использования

### Go TODO API (прямой доступ)

```bash
# Получить список задач
curl http://localhost:8080/todos

# Создать задачу
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "Изучить Go"}'

# Обновить задачу
curl -X PUT http://localhost:8080/todos/{id} \
  -H "Content-Type: application/json" \
  -d '{"title": "Изучить Go", "completed": true}'

# Удалить задачу
curl -X DELETE http://localhost:8080/todos/{id}
```

### FastAPI (с валидацией и приоритетами)

```bash
# Создать задачу с приоритетом
curl -X POST http://localhost:8000/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "Задача", "priority": "high"}'

# Получить статистику
curl http://localhost:8000/stats
```

### API Gateway

```bash
# Через шлюз к Go-сервису
curl http://localhost:9000/api/go/todos

# Через шлюз к FastAPI
curl http://localhost:9000/api/python/todos

# Проверка здоровья
curl http://localhost:9000/health
```

### WebSocket чат

```bash
cd python-service && . venv/bin/activate
python3 ws_client.py --username Alice --message "Привет всем!"
```

---

## Статический анализ

**Go:**
```
cd go-service && go vet ./...
# (нет ошибок)
```

**Python:**
```
cd python-service && . venv/bin/activate && ruff check app/ tests/ ws_client.py
# All checks passed!
```

---

## Нагрузочное тестирование

Смотри [BENCHMARK.md](BENCHMARK.md) для результатов сравнения производительности FastAPI vs Gin.

**Краткий итог:**

| Метрика | Go (Gin) | Python (FastAPI) |
|---------|----------|-----------------|
| Requests/sec | 96 044 | 33 |
| Latency avg | 1.33 ms | 2 800 ms |

---

## Тестирование

### Go тесты

```bash
cd go-service
go test -race -v ./...
```

**Результат:** 14 тестов, все PASS
- `store`: CRUD операции, thread-safety
- `handlers/http`: CRUD endpoints, валидация
- `middleware`: структурированное логирование
- `gateway`: reverse proxy, маршрутизация
- `ws`: Hub broadcast, disconnect handling

### Python тесты

```bash
cd python-service
. venv/bin/activate
python3 -m pytest -v
```

**Результат:** 13 тестов, все PASS
- `test_todos.py`: CRUD endpoints, валидация, обработка ошибок
- `test_stats.py`: статистика, пустой список
- `test_ws_client.py`: WebSocket клиент, JSON формат

---

## История коммитов

```
8b8f0a9 test(python): add WebSocket client tests
f83ff44 feat(python): add WebSocket client CLI script
c2af45d feat(ws): add WebSocket chat server with Hub pattern and tests
1c9fef2 chore: add gateway service to docker-compose and Dockerfile
ca3b327 feat(gateway): add API gateway with path-based routing and tests
42100ee docs(bench): add BENCHMARK.md with load test results
535d9e0 feat(bench): add wrk benchmark script
b304146 chore: add docker-compose.yml with service orchestration
f9cc4c6 chore: add Dockerfiles for Go and Python services
a2b0b43 feat(python): add stats endpoint and service unavailability handling with tests
3d36b24 feat(python): add FastAPI CRUD proxy endpoints with tests
a66e1a7 feat(python): add Pydantic models and Go service client with tests
74ffb45 chore(python): init FastAPI project with venv and dependencies
9fba428 refactor(go): apply handler/http per-scope pattern with app bootstrap
a4fde3e refactor(go): reorganize to golang-standards project layout
31d7430 fix(go): use gin.New with explicit recovery to avoid duplicate logging
e5a8679 feat(go): add structured logging middleware with tests
8556de7 feat(go): add Gin CRUD handlers for TODO API with tests
bdcaf4e feat(go): add TODO model and thread-safe in-memory store with tests
d80cd0f chore: init monorepo structure with README and gitignore
```

**Всего:** 20 атомарных коммитов с осмысленной историей.
