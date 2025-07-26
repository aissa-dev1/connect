package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type MessageResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type DataResponse[T any] struct {
	Data    T    `json:"data"`
	Success bool `json:"success"`
}

func RespondMessage(ctx *gin.Context, message string) MessageResponse {
	respone := MessageResponse{ Message: message, Success: true }
	ctx.JSON(http.StatusOK, respone)
	return respone
}

func RespondData[T any](ctx *gin.Context, data T) DataResponse[T] {
	response := DataResponse[T]{ Data: data, Success: true }
	ctx.JSON(http.StatusOK, response)
	return response
}

func RespondError(ctx *gin.Context, status int, message string) MessageResponse {
	response := MessageResponse{ Message: message, Success: false }
	ctx.JSON(status, response)
	return response
}

func RespondInternalError(ctx *gin.Context, message string) MessageResponse {
	return RespondError(ctx, http.StatusInternalServerError, message)
}

func RespondBadRequestError(ctx *gin.Context, message string) MessageResponse {
	return RespondError(ctx, http.StatusBadRequest, message)
}

func RespondNotFoundError(ctx *gin.Context, message string) MessageResponse {
	return RespondError(ctx, http.StatusNotFound, message)
}

func RespondUnauthorizedError(ctx *gin.Context, message string) MessageResponse {
	return RespondError(ctx, http.StatusUnauthorized, message)
}

func RespondForbiddenError(ctx *gin.Context, message string) MessageResponse {
	return RespondError(ctx, http.StatusForbidden, message)
}

func RespondUnprocessableError(ctx *gin.Context, message string) MessageResponse {
	return RespondError(ctx, http.StatusUnprocessableEntity, message)
}

func RespondConflictError(ctx *gin.Context, message string) MessageResponse {
	return RespondError(ctx, http.StatusConflict, message)
}
