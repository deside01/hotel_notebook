# Hotel Booking

Веб-приложение для бронирования отелей. Репозиторий содержит Go-бэкенд и Vue 3-фронтенд.

## Технологии

- **Backend:** Go, Chi Router
- **Frontend:** Vue 3, Vite
- **Database:** PostgreSQL 16
- **Migrations:** Goose
- **Infrastructure:** Docker, Docker Compose
- **Authentication:** JWT в `HttpOnly` cookies

## Структура проекта

```text
.
├── backend/
│   ├── cmd/
│   └── database/
│       └── migrations/
├── frontend/
├── docker-compose.yml
└── README.md
```

## Требования

- Go 1.21+
- Node.js и npm
- Docker и Docker Compose

## Запуск проекта

### 1. Запуск базы данных и миграций

Из корня репозитория выполните:

```bash
docker compose up -d
```

Команда запускает PostgreSQL и применяет миграции из `backend/database/migrations`, если этот путь указан в конфигурации Docker Compose.

### 2. Запуск бэкенда

Перейдите в директорию `backend`:

```bash
cd backend
go build -o backend.exe ./cmd
```

Запустите приложение в Windows PowerShell:

```powershell
.\backend.exe
```

По умолчанию API доступен по адресу `http://localhost:8080`, если порт не изменён в конфигурации.

### 3. Запуск фронтенда

Откройте второй терминал и перейдите в директорию `frontend`:

```bash
cd frontend
npm install
```

Создайте файл `.env.local` на основе `.env.example` и укажите адрес API:

```env
VITE_API_URL=http://localhost:8080
```

Запустите сервер разработки:

```bash
npm run dev
```

Фронтенд будет доступен по адресу `http://localhost:5173`.

### 4. Сборка фронтенда

Из директории `frontend` выполните:

```bash
npm run build
```

Для локального просмотра production-сборки:

```bash
npm run preview
```
