# URL Shortener

A robust URL shortener service built with Go, Gin, PostgreSQL, and Redis. This project provides a scalable API for shortening URLs, tracking usage, and managing links securely.

## Features

- **Shorten URLs:** Generate short links for long URLs.
- **Redirects:** Automatically redirect short URLs to their original destinations.
- **Analytics:** Track usage statistics for each short URL.
- **Validation:** Custom validators for URL formats and input data.
- **Persistence:** Store URL mappings in PostgreSQL.
- **Caching:** Use Redis for fast lookups and rate limiting.
- **Configurable:** Environment-based configuration for easy deployment.
- **Logging:** Structured logging for debugging and monitoring.
- **Secure:** Trusted proxies support for correct client IP handling.

## Project Structure

```
url-shortner/
├── main.go
├── config/
├── db/
├── routes/
├── utils/
├── validator/
└── readme.md
```

### Main Components

- **main.go:** Entry point. Initializes all services and starts the Gin server.
- **config:** Loads environment variables and app configuration.
- **db:** Handles connections to PostgreSQL and Redis.
- **routes:** Defines API endpoints and request handlers.
- **utils:** Utility functions, including logging.
- **validator:** Custom input validators for request data.

## Flow Overview

1. **Startup:**
  - Logger initialized (`utils.InitLogger`)
  - Configuration loaded from environment (`config.LoadConfig`)
  - Redis and PostgreSQL clients initialized (`db.InitRedis`, `db.InitDB`)
  - Custom validators registered (`validator.LoadCustomBindings`)
  - API routes set up (`routes.SetUpRouter`)
  - Trusted proxies configured for security
  - Gin server started

2. **Shorten URL:**
  - User sends a POST request with a long URL.
  - Input validated using custom validators.
  - Short URL generated and stored in PostgreSQL.
  - Mapping cached in Redis for fast access.

3. **Redirect:**
  - User accesses a short URL.
  - Service looks up the original URL in Redis (fallback to PostgreSQL).
  - Redirects user to the original URL.
  - Usage statistics updated.

## Configuration

Set environment variables for:

- `APP_HOST` and `APP_PORT`: Server address
- `TRUSTED_PROXIES`: Comma-separated list of trusted proxy IPs
- Database and Redis connection details

## Running Locally

```bash
go mod tidy
go run main.go
```

## API Endpoints

- `POST /shorten`: Shorten a new URL
- `GET /:shortUrl`: Redirect to the original URL
- `GET /stats/:shortUrl`: Get analytics for a short URL

## Dependencies

- [Gin](https://github.com/gin-gonic/gin): HTTP web framework
- [Redis](https://github.com/go-redis/redis): Caching and rate limiting
- [PostgreSQL](https://github.com/lib/pq): Persistent storage
- Custom packages for config, validation, logging, and routing

## License

MIT

---

**Contributions welcome!**