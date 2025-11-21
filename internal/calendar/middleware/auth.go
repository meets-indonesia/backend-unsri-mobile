package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"unsri-backend/internal/shared/errors"
	"unsri-backend/internal/shared/utils"
	"unsri-backend/pkg/jwt"
)

// AuthMiddleware validates JWT token
func AuthMiddleware(jwtToken *jwt.JWT) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, 401, errors.NewUnauthorizedError("authorization header required"))
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(c, 401, errors.NewUnauthorizedError("invalid authorization header format"))
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := jwtToken.ValidateToken(token)
		if err != nil {
			utils.ErrorResponse(c, 401, errors.NewUnauthorizedError("invalid or expired token"))
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)
		c.Set("user_email", claims.Email)

		c.Next()
	}
}

// RoleMiddleware checks if user has required role
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			utils.ErrorResponse(c, 401, errors.NewUnauthorizedError("user role not found"))
			c.Abort()
			return
		}

		role := userRole.(string)
		allowed := false
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				allowed = true
				break
			}
		}

		if !allowed {
			utils.ErrorResponse(c, 403, errors.NewForbiddenError("insufficient permissions"))
			c.Abort()
			return
		}

		c.Next()
	}
}

