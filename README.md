# 🛡️ Auth Service — Golang микросервис авторизации

![Go CI](https://github.com/dmendixan/go_auth_service/actions/workflows/go.yml/badge.svg)


Полноценный микросервис на Go, реализующий систему аутентификации и авторизации с поддержкой JWT, refresh-токенов, ролевой модели и REST API. Используется в рамках производственной практики в АО "Казахтелеком".

---

## 🔧 Возможности

- 📦 Регистрация и логин пользователя
- 🔐 Генерация JWT и refresh токенов
- 👮 Ролевая модель: `user` и `admin`
- 🧾 Защищённые эндпоинты с middleware
- 🗂 CRUD-функции для администратора
- 🧪 Полные unit-тесты
- 🐳 Docker + PostgreSQL
- ⚙️ CI через GitHub Actions

---

## 🚀 Быстрый старт (через Docker Compose)

### 1. Клонируйте проект
```bash
git clone https://github.com/dmendixan/go_auth_service.git
cd auth-service
```

### 2. Создайте `.env` файл в корне:
```env
DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=123321
DB_NAME=auth_service
JWT_SECRET=supersecretkey
TOKEN_EXPIRY=1h
```

### 3. Запустите проект:
```bash
docker-compose up --build
```

📍 API доступен по адресу: [http://localhost:8080](http://localhost:8080)


---

## 🗃️ Структура проекта
```bash
auth-service/
├── cmd/
│   └── main.go                # Точка входа в приложение
├── config/                    # Инициализация базы данных, переменные окружения
├── internal/
│   ├── handlers/              # HTTP-обработчики
│   ├── models/                # GORM модели (User, RefreshToken)
│   ├── repository/            # Работа с БД (опционально)
│   └── services/              # JWT, токены, бизнес-логика
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── .env
```

---

## 🧪 Unit-тесты

```bash
go test ./...
```

✔️ Покрытие: регистрация, логин, refresh, защищённые эндпоинты, действия администратора.

Тестируются следующие сценарии:
- `/api/register`
- `/api/login`
- `/api/refresh`
- `/api/profile`
- `/api/admin/users`
- `/api/admin/users/delete`
- `/api/admin/users/update`
---

## 📬 Эндпоинты

| Метод | Endpoint              | Описание                                |
|-------|------------------------|-----------------------------------------|
| POST  | `/api/register`       | Регистрация пользователя                |
| POST  | `/api/login`          | Аутентификация и выдача токенов         |
| POST  | `/api/refresh`        | Обновление access токена                |
| GET   | `/api/profile`        | Данные профиля (требуется JWT)          |
| GET   | `/api/admin/users`    | Получить всех пользователей (admin)     |
| PUT   | `/api/admin/users/:id`| Обновить пользователя (admin)           |
| DELETE| `/api/admin/users/:id`| Удалить пользователя (admin)            |

---

## ⚙️ CI/CD через GitHub Actions

Файл: `.github/workflows/go.yml`

- При каждом push в `main`:
  - запускается PostgreSQL
  - выполняется `go test ./...`
  - проверяется работоспособность проекта

![Go CI](https://github.com/dmendixan/go_auth_service/actions/workflows/go.yml/badge.svg)

---

## 🧠 Автор

> **Меңдіхан Дінмұхаммед**  
> Практикант в АО "Казахтелеком"  
> Backend-разработка на Go  
> Telegram: [@dmendixan](https://t.me/dmendixan)



---

## ❗ Возможные ошибки

- **"No .env file found"** — убедитесь, что `.env` скопирован и передан в Docker context
- **"db: no such host"** — используйте `localhost` при локальном запуске без Docker
Чтобы проверить эндпоинт от админа, назначьте себя админом в бд

---

## ✅ Готово к проверке

- Микросервис запускается одной командой через Docker
- Все функции протестированы и задокументированы
- Код сопровождается тестами и CI-проверками

📦 Проект готов к использованию и масштабированию ✨
