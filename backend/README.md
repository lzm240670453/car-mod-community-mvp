# retrofit backend

Go backend skeleton for 「撸车日记」.

## Local setup

```powershell
Copy-Item .env.example .env
docker compose up -d
go mod tidy
go run ./cmd/api
```

Health check:

```text
GET http://localhost:8080/healthz
```
