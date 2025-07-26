package blockhandler

import (
	blockconstants "connect/internal/block/constants"
	blockmodel "connect/internal/block/model"
	blockservice "connect/internal/block/service"
	"connect/internal/db"
	friendshipservice "connect/internal/friendship/service"
	errormessage "connect/internal/pkg/error_message"
	"connect/internal/pkg/response"
	profileservice "connect/internal/profile/service"
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
)

func BlockHandler(ctx *gin.Context) {
	blockedIdParam := ctx.Param("blockedId")
	blockedId, blockedIdErr := strconv.Atoi(blockedIdParam)

	if blockedIdErr != nil {
		response.RespondInternalError(ctx, errormessage.InternalServerError)
		return
	}

	profileId, profileIdErr := profileservice.GetProfileId(ctx)

	if profileIdErr != nil {
		response.RespondInternalError(ctx, profileIdErr.Error())
		return
	}
	if profileId == blockedId {
		response.RespondForbiddenError(ctx, "You can't block yourself")
		return
	}

	blockExists, blockExistsErr := blockservice.BlockExists(profileId, blockedId)

	if blockExistsErr != nil {
		response.RespondInternalError(ctx, blockExistsErr.Error())
		return
	}
	if blockExists {
		response.RespondUnprocessableError(ctx, "You have already blocked this user")
		return
	}

	blockErr := blockservice.InsertBlock(blockmodel.Block{ BlockerId: profileId, BlockedId: blockedId })

	if blockErr != nil {
		response.RespondInternalError(ctx, blockErr.Error())
		return
	}

	friendshipExists, friendshipExistsErr := friendshipservice.FriendshipExists(profileId, blockedId)

	if friendshipExistsErr != nil {
		response.RespondInternalError(ctx, friendshipExistsErr.Error())
		return
	}
	if friendshipExists {
		friendshipservice.DeleteFriendship(profileId, blockedId)
	}

	response.RespondMessage(ctx, "User have been blocked successfully")
}

func UnblockHandler(ctx *gin.Context) {
	blockedIdParam := ctx.Param("blockedId")
	blockedId, blockedIdErr := strconv.Atoi(blockedIdParam)

	if blockedIdErr != nil {
		response.RespondInternalError(ctx, errormessage.InternalServerError)
		return
	}

	profileId, profileIdErr := profileservice.GetProfileId(ctx)

	if profileIdErr != nil {
		response.RespondInternalError(ctx, profileIdErr.Error())
		return
	}

	blockExists, blockExistsErr := blockservice.BlockExists(profileId, blockedId)

	if blockExistsErr != nil {
		response.RespondInternalError(ctx, blockExistsErr.Error())
		return
	}
	if !blockExists {
		response.RespondNotFoundError(ctx, blockconstants.ErrBlockNotFound.Error())
		return
	}

	_, unblockErr := db.Pool().Exec(context.Background(), "DELETE FROM blocks WHERE blockerId = $1 AND blockedId = $2;", profileId, blockedId)

	if unblockErr != nil {
		response.RespondInternalError(ctx, errormessage.InternalServerError)
		return
	}

	response.RespondMessage(ctx, "User have been unblocked successfully")
}
