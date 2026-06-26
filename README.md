# Ticket System API

A simple backend service written in Go for a ticket system. It supports user registration, login, ticket creation, ownership-based ticket isolation, and custom ticket status validation rules.

## Deployed URL
* **Live Base URL**: `https://ticket-system-sybl.onrender.com`
* **Interactive API Documentation (Swagger)**: [https://ticket-system-sybl.onrender.com/docs](https://ticket-system-sybl.onrender.com/docs)
* **Health Check**: [https://ticket-system-sybl.onrender.com/health](https://ticket-system-sybl.onrender.com/health)

## Tech Stack
- Go
- PostgreSQL
- Chi (Router)
- JWT (for auth)
- Bcrypt (for password hashing)

## Project Setup

Create a `.env` file in the project root:
```env
PORT=8080
JWT_SECRET=some_jwt_secret_key
DATABASE_URL=postgres://postgres:password@localhost:5432/postgres?sslmode=disable
```

## Running Locally

To run the application directly, make sure you have PostgreSQL running, then execute:
```bash
go run cmd/server/main.go
```
The application automatically creates the required tables (`users`, `tickets`, `user_tickets`) on startup if they don't exist.

## Running in Docker

To build and run using Docker:

```bash
# Build the image
docker build -t ticket-system .

# Run the container
docker run -p 8080:8080 -e DATABASE_URL="postgres://postgres:tiger@host.docker.internal:5432/postgres?sslmode=disable" ticket-system
```

## API Endpoints

Once the server is running, you can access the interactive Swagger documentation at:
**http://localhost:8080/docs**

### Endpoints:
- `GET /health` - Health check (returns `{"status":"ok"}`)
- `POST /auth/register` - Create an account
- `POST /auth/login` - Authenticate and get JWT token
- `POST /tickets` - Create a ticket (authenticated)
- `GET /tickets` - List own tickets (authenticated)
- `GET /tickets/{id}` - Fetch details for own ticket (authenticated)
- `PATCH /tickets/{id}/status` - Update status of own ticket (authenticated)

## Deployment on Render

1. Create a free **PostgreSQL** database on Render. Copy the Internal Database URL.
2. Create a new **Web Service** on Render connected to your GitHub repo.
3. Select **Docker** as the runtime environment.
4. Add environment variables:
   - `DATABASE_URL` (Internal database URL from step 1)
   - `JWT_SECRET` (Your JWT signing secret)
