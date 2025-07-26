package blockhandler

import (
	blockmodel "connect/internal/block/model"
	blockservice "connect/internal/block/service"
	errormessage "connect/internal/pkg/error_message"
	"connect/internal/pkg/response"
	profileservice "connect/internal/profile/service"
	usermodel "connect/internal/user/model"
	userservice "connect/internal/user/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetBlocksBetweenUsersHandler(ctx *gin.Context) {
	blockerIdParam := ctx.Param("blockerId")
	blockerId, blockerIdErr := strconv.Atoi(blockerIdParam)

	if blockerIdErr != nil {
		response.RespondInternalError(ctx, errormessage.InternalServerError)
		return
	}

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
	if profileId != blockerId && profileId != blockedId {
		response.RespondForbiddenError(ctx, errormessage.ForbiddenResource)
		return
	}

	blocks := []blockmodel.Block{}

	block1, block1Err := blockservice.GetBlock(blockerId, blockedId)

	if block1Err != nil {
		response.RespondInternalError(ctx, block1Err.Error())
		return
	}

	block2, block2Err := blockservice.GetBlock(blockedId, blockerId)

	if block2Err != nil {
		response.RespondInternalError(ctx, block2Err.Error())
		return
	}

	if block1 != nil {
		blocks = append(blocks, *block1)
	}
	if block2 != nil {
		blocks = append(blocks, *block2)
	}

	response.RespondData[[]blockmodel.Block](ctx, blocks)
}

func GetBlockedUsersHandler(ctx *gin.Context) {
	profileId, profileIdErr := profileservice.GetProfileId(ctx)

	if profileIdErr != nil {
		response.RespondInternalError(ctx, profileIdErr.Error())
		return
	}

	blocks, blocksErr := blockservice.GetBlockerBlocks(profileId)

	if blocksErr != nil {
		response.RespondInternalError(ctx, blocksErr.Error())
		return
	}

	users := []usermodel.MinimalUser{}

	for _, block := range blocks {
		user, userErr := userservice.GetMinimalUserById(block.BlockedId)

		if userErr != nil {
			response.RespondInternalError(ctx, userErr.Error())
			return
		}
		if user == nil {
			continue
		}

		users = append(users, *user)
	}

	response.RespondData[[]usermodel.MinimalUser](ctx, users)
}
