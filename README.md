# TaskSphere

Full-stack task manager built with Go (Fiber), PostgreSQL, and React/Vite.

## Backend

1. Copy `.env.example` (if needed) and set `DATABASE_URL`, `JWT_SECRET`, `PORT`.
2. Run PostgreSQL + Adminer + backend with Docker:
   ```bash
   docker compose up --build
   ```
   - Backend: `http://localhost:8080` (if running inside Docker)
   - Adminer: `http://localhost:8082` (system: PostgreSQL, server: `db`, user `admin`, password `secret`, database `tasksphere`)
3. Or run backend locally:
   ```bash
   cd backend
   go run main.go
   ```

## Frontend

1. Install deps and start dev server (port auto-selects 5173-5200):
   ```bash
   npm --prefix frontend install
   ./scripts/dev.sh
   ```
2. Configure API base via `frontend/.env` (`VITE_API_URL=http://localhost:8081`, etc.).

## SQLC

SQL queries live in `backend/db/queries`. Generated Go code is committed in `backend/db/sqlc`.

- Update queries/schema → regenerate:
  ```bash
  sqlc generate
  ```
- Config: `sqlc.yaml` (PostgreSQL engine, outputs Go code compatible with pgx/v5).
