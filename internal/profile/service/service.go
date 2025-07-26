package profileservice

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetProfileId(ctx *gin.Context) (int, error) {
	val, exists := ctx.Get("profileId")

	if !exists {
		return 0, errors.New("profileId not found in context")
	}

	id, ok := val.(int)

	if !ok {
		return 0, fmt.Errorf("profileId has unexpected type: %s", val)
	}

	return id, nil
}
