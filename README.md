# linkcheck

Конкурентный CLI-чекер ссылок на Go: параллельная проверка списка URL с
ограничением нагрузки, отменой через context и graceful shutdown.

## Структура проекта

```
linkchecker/
├── cmd/
│   └── linkcheck/        точка входа CLI, флаги
├── internal/
│   ├── checker/          Result, worker pool, семафор/rate limiter, context, select (Этапы 1-5)
│   ├── cache/            три реализации кэша: Mutex / RWMutex / sync.Map + бенчмарк (Этап 6)
│   ├── pipeline/         конвейер стадий gen→normalize→dedup→check→collect (Этап 10)
│   ├── server/           HTTP-сервер с /check и graceful shutdown (Этап 9)
│   └── testserver/       локальный httptest-сервер для отладки без сети
└── reports/              JSON-отчёты прогонов (в git не попадают)
```

## Флаги CLI

| Флаг       | Тип      | Описание                                   |
|------------|----------|---------------------------------------------|
| `-workers` | int      | размер пула воркеров (по умолчанию NumCPU)  |
| `-rps`     | int      | запросов в секунду                          |
| `-timeout` | duration | общий дедлайн на весь прогон                |
| `-per-req` | duration | таймаут одного запроса                      |
| `-out`     | string   | файл отчёта (JSON)                          |

## Запуск

```bash
go run ./cmd/linkcheck -workers 8 -rps 10 -timeout 30s -per-req 3s -out reports/report.json https://example.com https://golang.org
```

## Разработка и отладка

Локальный тестовый сервер (`internal/testserver`) поднимает ручки:

- `/status/{code}` — сразу отвечает указанным кодом;
- `/delay/{ms}` — отвечает 200 после задержки в мс;
- `/timeout` — никогда не отвечает (для проверки таймаута запроса).

## Бенчмарк кэша (Этап 6)

```bash
go test ./internal/cache/... -bench=. -benchmem
```

<!-- TODO: сюда результаты бенчмарка и вывод про sync.Map vs RWMutex -->

## Проверки

```bash
go vet ./...
go test -race ./...
go run -race ./cmd/linkcheck ...
```

## Наблюдение за планировщиком и GC (Этап 11)

```bash
GOMAXPROCS=1 go run ./cmd/linkcheck ...
GODEBUG=schedtrace=1000 go run ./cmd/linkcheck ...
GODEBUG=gctrace=1 go run ./cmd/linkcheck ...
```
