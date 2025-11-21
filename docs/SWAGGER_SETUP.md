# Swagger Setup Guide

Panduan untuk setup Swagger documentation.

## Installation

Swagger dependencies sudah terinstall. Untuk generate documentation:

```bash
# Install swag CLI (if not installed)
go install github.com/swaggo/swag/cmd/swag@latest

# Generate swagger docs
swag init -g cmd/api-gateway/main.go
```

Ini akan membuat folder `docs/` dengan file:
- `docs/swagger.json`
- `docs/swagger.yaml`
- `docs/docs.go`

## Enable Swagger UI

1. Set environment variable:
```bash
export ENABLE_SWAGGER=true
```

2. Run API Gateway:
```bash
make run-api-gateway
```

3. Akses Swagger UI:
```
http://localhost:8080/swagger/index.html
```

## Adding Swagger Annotations

Tambahkan annotations di handler functions:

```go
// @Summary Get user profile
// @Description Get current user profile
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User
// @Failure 401 {object} utils.Response
// @Router /users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
    // ...
}
```

## Swagger Annotations Reference

- `@Summary` - Short summary
- `@Description` - Detailed description
- `@Tags` - Group endpoints
- `@Accept` - Content types accepted
- `@Produce` - Content types produced
- `@Security` - Security scheme
- `@Param` - Parameter description
- `@Success` - Success response
- `@Failure` - Error response
- `@Router` - Route path and method

Lihat [Swaggo Documentation](https://github.com/swaggo/swag) untuk lebih detail.

