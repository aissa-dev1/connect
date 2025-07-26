package authservice

import (
	jwtpkg "connect/internal/pkg/jwt"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterToken(ctx *gin.Context, id int) (string, error) {
	token, jwtGenErr := jwtpkg.GenerateJWT(fmt.Sprintf("%d", id))

	if jwtGenErr != nil {
		return "", errors.New("cannot generate jwt token from sub")
	}

	domain := ""
	isProd := os.Getenv("ENVIRONMENT") == "production"

	if isProd {
		domain = os.Getenv("CLIENT_DOMAIN")
	}

	ctx.SetCookie("token", token, int(time.Until(jwtpkg.GetJWTExpirationTime()).Seconds()), "/", domain, isProd, true)

	return token, nil
}

func ClearToken(ctx *gin.Context) {
	domain := ""
	isProd := os.Getenv("ENVIRONMENT") == "production"

	if isProd {
		domain = os.Getenv("CLIENT_DOMAIN")
	}

	ctx.SetCookie("token", "", -1, "/", domain, isProd, true)
}
