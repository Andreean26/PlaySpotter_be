package middlewares

import (
	"net/http"

	"playspotter/internal/utils"

	"github.com/gin-gonic/gin"
)

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := GetUserRole(c)
		if !exists {
			utils.RespondError(c, http.StatusForbidden, "forbidden", "User role not found")
			c.Abort()
			return
		}

		// Check if user's role is in the allowed roles
		allowed := false
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				allowed = true
				break
			}
		}

		if !allowed {
			utils.RespondError(c, http.StatusForbidden, "forbidden", "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}
