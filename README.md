# **Сервис назначения ревьюеров для Pull Request (Осень 2025)**

Мини-сервис на Go, который автоматически назначает ревьюеров на PR, позволяет переназначать одного из ревьюеров, управлять командами и активностью пользователей. Спецификация — в `task/openapi.yml`.

## Запуск
- Требования: Docker, Docker Compose.
- Порт: `8080`.

```bash
# основной запуск
docker compose -f deploy/docker-compose.yml up -d

# проверка живости и готовности
curl -s http://localhost:8080/health
curl -s http://localhost:8080/ready
```

Миграции применяются автоматически через сервис `migrator`.

## Основные эндпоинты (curl)

Команды без переменных окружения, можно копировать как есть.

Здоровье:
```bash
curl -i http://localhost:8080/health
```

Команды/пользователи:
```bash
# создать команду
curl -s -X POST http://localhost:8080/team/add \
  -H 'Content-Type: application/json' \
  -d '{"team_name":"team-a","members":[
    {"user_id":"u10","username":"Neo","is_active":true},
    {"user_id":"u11","username":"Trin","is_active":true},
    {"user_id":"u12","username":"Morp","is_active":false}
  ]}'

# повторное создание -> TEAM_EXISTS
curl -s -X POST http://localhost:8080/team/add \
  -H 'Content-Type: application/json' \
  -d '{"team_name":"team-a","members":[{"user_id":"x","username":"dup","is_active":true}]}'

# получить команду
curl -s "http://localhost:8080/team/get?team_name=team-a"

# деактивация пользователя
curl -s -X POST http://localhost:8080/users/setIsActive \
  -H 'Content-Type: application/json' \
  -d '{"user_id":"u11","is_active":false}'

# PR'ы, где пользователь ревьювер
curl -s "http://localhost:8080/users/getReview?user_id=u11"
```

PR: создание, merge (идемпотентно), переназначение:
```bash
# создать PR (назначаются до 2 активных из команды автора, без автора)
curl -s -X POST http://localhost:8080/pullRequest/create \
  -H 'Content-Type: application/json' \
  -d '{"pull_request_id":"pr-a-1","pull_request_name":"Init A1","author_id":"u10"}'

# merge (повторный вызов возвращает актуальное MERGED состояние)
curl -s -X POST http://localhost:8080/pullRequest/merge \
  -H 'Content-Type: application/json' \
  -d '{"pull_request_id":"pr-a-1"}'

# переназначение одного из назначенных ревьюверов
curl -s -X POST http://localhost:8080/pullRequest/reassign \
  -H 'Content-Type: application/json' \
  -d '{"pull_request_id":"pr-a-1","old_user_id":"<реально_назначенный_id>"}'
```

Статистика:
```bash
curl -s http://localhost:8080/stats | jq
```

## E2E-проверка (docker-compose)

Файл: `deploy/docker-compose.e2e.yml`
```bash
docker compose -f deploy/docker-compose.e2e.yml up --build tester --exit-code-from tester
```
Сценарий: миграции -> поднятие API -> ожидание `/ready` -> создание команды -> создание PR -> merge -> проверка статуса.

## Нагрузочное тестирование (k6)
Простой скрипт в `scripts/k6/pr_flow.js`:
```bash
k6 run scripts/k6/pr_flow.js
```
Сценарий создает команду, PR, делает merge. Цель: RPS≈5, p95<300мс на локальной машине.

## Линтер и тесты
```bash
# линтер
golangci-lint run

# юнит-тесты
go test ./...
```

## Замечания по реализации
- Назначение ревьюверов: до 2 активных из команды автора, автор исключается.
- Переназначение: из команды «старого» ревьювера, активный, не автор, не уже назначенный.
- После MERGED изменения запрещены (PR_MERGED).
- merge — идемпотентный.
- Возвращаемые поля PR включают `assigned_reviewers`, `createdAt`, `mergedAt`.
