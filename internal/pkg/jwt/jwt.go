package jwtpkg

import (
	"connect/internal/pkg/response"
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	Sub string `json:"id"`
	jwt.RegisteredClaims
}

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

func GenerateJWT(sub string) (string, error) {
	claims := JWTClaims{
		Sub: sub,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(GetJWTExpirationTime()),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Issuer: "connect",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtKey)
}

func ValidateJWT(tokenStr string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)

	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func GetJWTExpirationTime() time.Time {
	expHours, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION_HOURS"))
	
	if err != nil {
		log.Printf("Failed to convert jwt expiration hours %v\n", err)
		return time.Now().Add(0 * time.Hour)
	}

	return time.Now().Add(time.Duration(expHours) * time.Hour)
}

func ApplyJWTGuard() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, tokenErr := ctx.Cookie("token")

		if tokenErr != nil {
			response.RespondUnauthorizedError(ctx, "Unauthorized: No token")
			ctx.Abort()
			return 
		}

		claims, claimsErr := ValidateJWT(token)

		if claimsErr != nil {
			response.RespondUnauthorizedError(ctx, "Unauthorized: Invalid token")
			ctx.Abort()
			return 
		}

		profileId, profileIdErr := strconv.Atoi(claims.Sub)

		if profileIdErr != nil {
			response.RespondInternalError(ctx, "Failed to convert profile id from a string to int")
			ctx.Abort()
			return 
		}

		ctx.Set("profileId", profileId)
		ctx.Next()
	}
}
