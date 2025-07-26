package profilehandler

import (
	"connect/internal/db"
	errormessage "connect/internal/pkg/error_message"
	"connect/internal/pkg/response"
	profilemodel "connect/internal/profile/model"
	profileservice "connect/internal/profile/service"
	"context"

	"github.com/gin-gonic/gin"
)

func GetMyProfileHandler(ctx *gin.Context) {
	profileId, profileIdErr := profileservice.GetProfileId(ctx)

	if profileIdErr != nil {
		response.RespondInternalError(ctx, profileIdErr.Error())
		return
	}

	var profile profilemodel.Profile

	getUserRowErr := db.Pool().QueryRow(context.Background(), "SELECT id, email, username FROM users WHERE id = $1;", profileId).Scan(&profile.Id, &profile.Email, &profile.Username)

	if getUserRowErr != nil {
		response.RespondInternalError(ctx, errormessage.InternalServerError)
		return
	}

	response.RespondData[profilemodel.Profile](ctx, profile)
}
