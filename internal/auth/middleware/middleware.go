package authmiddleware

import (
	jwtpkg "connect/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return jwtpkg.ApplyJWTGuard()
}
