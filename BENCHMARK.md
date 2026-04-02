# Нагрузочное тестирование — FastAPI vs Gin

## Методология

- **Инструмент**: wrk (debian/4.1.0-4build2 [epoll])
- **Параметры**: `-t4 -c100 -d30s --latency`
- **Эндпоинт**: `GET /todos` (оба сервиса)
- **Предварительная подготовка**: 10 TODO-записей засеяно в Go-сервис перед запуском тестов
- **Окружение**: Linux 6.17.0 x86_64, Intel Core Ultra 9 185H, 22 CPU cores

## Результаты

### Go (Gin) — порт 8080

```
Running 30s test @ http://localhost:8080/todos
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.33ms    1.45ms  26.13ms   86.07%
    Req/Sec    24.15k     1.84k   30.57k    72.58%
  Latency Distribution
     50%    1.14ms
     75%    2.05ms
     90%    3.20ms
     99%    6.09ms
  2883522 requests in 30.02s, 2.78GB read
Requests/sec:  96044.69
Transfer/sec:     94.89MB
```

### Python (FastAPI) — порт 8000

> FastAPI выступает прокси: каждый запрос `GET /todos` проксируется к Go-сервису через `httpx.AsyncClient`.
> Запуск с `--timeout 10s` для учёта задержки проксирования.

```
Running 30s test @ http://localhost:8000/todos
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     2.80s   890.57ms   9.97s    89.72%
    Req/Sec    15.97     25.46   171.00     89.89%
  Latency Distribution
     50%    2.71s
     75%    2.79s
     90%    2.91s
     99%    7.75s
  997 requests in 30.05s, 1.31MB read
  Socket errors: connect 0, read 0, write 0, timeout 11
Requests/sec:     33.18
Transfer/sec:     44.68KB
```

## Сравнение

| Метрика          | Go (Gin)      | Python (FastAPI) |
|------------------|---------------|-----------------|
| Requests/sec     | 96 044        | 33              |
| Latency avg      | 1.33 ms       | 2 800 ms        |
| Latency p50      | 1.14 ms       | 2 710 ms        |
| Latency p99      | 6.09 ms       | 7 750 ms        |
| Ошибки/таймауты  | 0             | 11              |
| Данных прочитано | 2.78 GB       | 1.31 MB         |

## Анализ

**Go (Gin) значительно превосходит Python (FastAPI) под нагрузкой**:

- **Throughput**: Gin обрабатывает ~2 900× больше запросов в секунду (96 044 vs 33 req/sec)
- **Latency**: средняя задержка Gin ниже в ~2 100× (1.33ms vs 2 800ms)
- **Стабильность**: у Gin 0 ошибок против 11 таймаутов у FastAPI

**Причины разрыва**:

1. **Архитектура**: FastAPI является прокси-сервисом — каждый запрос `GET /todos` делает дополнительный HTTP-вызов к Go-сервису через `httpx.AsyncClient`. Это удваивает сетевые задержки и добавляет накладные расходы на создание HTTP-соединения при каждом запросе.

2. **Компилируемый vs интерпретируемый**: Go компилируется в нативный бинарник без GIL. FastAPI работает на CPython с async event loop (uvicorn + asyncio), что эффективно для I/O-bound задач, но имеет более высокий накладной расход на запрос.

3. **Хранилище**: Go обращается к in-memory store напрямую через `sync.RWMutex` — это атомарная операция в памяти. FastAPI, помимо собственной логики, делает сетевой запрос к Go.

4. **Создание клиента**: В текущей реализации FastAPI создаёт новый `httpx.AsyncClient` на каждый запрос (`async with httpx.AsyncClient()`), что добавляет overhead на setup/teardown TCP-соединения.

**Вывод**: Разрыв в производительности обусловлен не только разницей в языках, но и принципиально разной архитектурой сервисов (прямой доступ vs прокси). При сравнении «чистого» FastAPI без проксирования разница была бы значительно меньше (обычно 3–10×).
