package messagehandler

import (
	messagemodel "connect/internal/message/model"
	messageservice "connect/internal/message/service"
	errormessage "connect/internal/pkg/error_message"
	"connect/internal/pkg/response"
	profileservice "connect/internal/profile/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetMessagesBetweenUsersHandler(ctx *gin.Context) {
	senderIdParam := ctx.Param("senderId")
	senderId, senderIdErr := strconv.Atoi(senderIdParam)

	if senderIdErr != nil {
		response.RespondInternalError(ctx, errormessage.InternalServerError)
		return
	}

	receiverIdParam := ctx.Param("receiverId")
	receiverId, receiverIdErr := strconv.Atoi(receiverIdParam)

	if receiverIdErr != nil {
		response.RespondInternalError(ctx, errormessage.InternalServerError)
		return
	}

	profileId, profileIdErr := profileservice.GetProfileId(ctx)

	if profileIdErr != nil {
		response.RespondInternalError(ctx, profileIdErr.Error())
		return
	}
	if profileId != senderId && profileId != receiverId {
		response.RespondForbiddenError(ctx, errormessage.ForbiddenResource)
		return
	}

	messages, messagesErr := messageservice.GetMessagesBetweenUsers(senderId, receiverId)

	if messagesErr != nil {
		response.RespondInternalError(ctx, messagesErr.Error())
		return
	}

	response.RespondData[[]messagemodel.Message](ctx, messages)
}
