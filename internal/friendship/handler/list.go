package friendshiphandler

import (
	friendshipconstants "connect/internal/friendship/constants"
	friendshipmodel "connect/internal/friendship/model"
	friendshipservice "connect/internal/friendship/service"
	errormessage "connect/internal/pkg/error_message"
	"connect/internal/pkg/response"
	profileservice "connect/internal/profile/service"
	usermodel "connect/internal/user/model"
	userservice "connect/internal/user/service"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetFriendRequestHandler(ctx *gin.Context) {
	receiverIdParam := ctx.Param("id")
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

	friendship, friendshipErr := friendshipservice.MustGetFriendship(profileId, receiverId)

	if friendshipErr != nil {
		if errors.Is(friendshipErr, friendshipconstants.ErrFriendshipNotFound) {
			response.RespondNotFoundError(ctx, friendshipErr.Error())
		} else {
			response.RespondInternalError(ctx, friendshipErr.Error())
		}
		return
	}

	response.RespondData[friendshipmodel.Friendship](ctx, friendship)
}

func GetAllFriendRequestsHandler(ctx *gin.Context) {
	receiverIdParam := ctx.Param("id")
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
	if profileId != receiverId {
		response.RespondForbiddenError(ctx, errormessage.ForbiddenResource)
		return
	}

	friendships, friendshipsErr := friendshipservice.GetReceiverFriendships(receiverId)

	if friendshipsErr != nil {
		response.RespondInternalError(ctx, friendshipsErr.Error())
		return
	}

	response.RespondData[[]friendshipmodel.Friendship](ctx, friendships)
}

func GetRequestsHandler(ctx *gin.Context) {
	profileId, profileIdErr := profileservice.GetProfileId(ctx)

	if profileIdErr != nil {
		response.RespondInternalError(ctx, profileIdErr.Error())
		return
	}

	friendships, friendshipsErr := friendshipservice.GetReceiverFriendshipsByStatus(profileId, friendshipconstants.StatusPending)

	if friendshipsErr != nil {
		response.RespondInternalError(ctx, friendshipsErr.Error())
		return
	}

	users := []usermodel.MinimalUser{}

	for _, friendship := range friendships {
		user, userErr := userservice.GetMinimalUserById(friendship.RequesterId)

		if userErr != nil {
			response.RespondInternalError(ctx, userErr.Error())
			break
		}
		if user == nil {
			continue
		}

		users = append(users, *user)
	}

	response.RespondData[[]usermodel.MinimalUser](ctx, users)
}

func GetFriendsHandler(ctx *gin.Context) {
	profileId, profileIdErr := profileservice.GetProfileId(ctx)

	if profileIdErr != nil {
		response.RespondInternalError(ctx, profileIdErr.Error())
		return
	}

	friendships, friendshipsErr := friendshipservice.GetRequesterOrReceiverFriendshipsByStatus(profileId, friendshipconstants.StatusAccepted)

	if friendshipsErr != nil {
		response.RespondInternalError(ctx, friendshipsErr.Error())
		return
	}

	users := []usermodel.MinimalUser{}

	for _, friendship := range friendships {
		var id int

		if profileId == friendship.RequesterId {
			id = friendship.ReceiverId
		} else {
			id = friendship.RequesterId
		}

		user, userErr := userservice.GetMinimalUserById(id)

		if userErr != nil {
			response.RespondInternalError(ctx, errormessage.InternalServerError)
			break
		}
		if user == nil {
			continue
		}

		users = append(users, *user)
	}

	response.RespondData(ctx, users)
}
