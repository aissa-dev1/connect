package friendshiphandler

import (
	blockservice "connect/internal/block/service"
	"connect/internal/db"
	friendshipconstants "connect/internal/friendship/constants"
	friendshipmodel "connect/internal/friendship/model"
	friendshipservice "connect/internal/friendship/service"
	errormessage "connect/internal/pkg/error_message"
	"connect/internal/pkg/response"
	profileservice "connect/internal/profile/service"
	"context"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SendFriendRequestHandler(ctx *gin.Context) {
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
	if profileId == receiverId {
		response.RespondUnprocessableError(ctx, "You can't send a friend request to yourself")
		return
	}

	blockExists, blockExistsErr := blockservice.BlockExistsMutual(profileId, receiverId)

	if blockExistsErr != nil {
		response.RespondInternalError(ctx, blockExistsErr.Error())
		return
	}
	if blockExists {
		response.RespondForbiddenError(ctx, "One of you blocked another")
		return
	}

	requestExists, requestExistsErr := friendshipservice.FriendshipExists(profileId, receiverId)

	if requestExistsErr != nil {
		response.RespondInternalError(ctx, requestExistsErr.Error())
		return
	}
	if requestExists {
		response.RespondConflictError(ctx, "One of you already sent a request to the other")
		return
	}

	insertFriendshipErr := friendshipservice.InsertFriendship(friendshipmodel.Friendship{RequesterId: profileId, ReceiverId: receiverId})

	if insertFriendshipErr != nil {
		response.RespondInternalError(ctx, errormessage.InternalServerError)
		return
	}

	response.RespondMessage(ctx, "Friend request sent successfully")
}

func AcceptFriendRequestHandler(ctx *gin.Context) {
	requesterIdParam := ctx.Param("requesterId")
	requesterId, requesterIdErr := strconv.Atoi(requesterIdParam)

	if requesterIdErr != nil {
		response.RespondInternalError(ctx, errormessage.InternalServerError)
		return
	}

	profileId, profileIdErr := profileservice.GetProfileId(ctx)

	if profileIdErr != nil {
		response.RespondInternalError(ctx, profileIdErr.Error())
		return
	}

	friendship, friendshipErr := friendshipservice.MustGetFriendship(requesterId, profileId)

	if friendshipErr != nil {
		if errors.Is(friendshipErr, friendshipconstants.ErrFriendshipNotFound) {
			response.RespondNotFoundError(ctx, friendshipErr.Error())
		} else {
			response.RespondInternalError(ctx, friendshipErr.Error())
		}
		return
	}
	if friendship.Status == nil || *friendship.Status != friendshipconstants.StatusPending {
		response.RespondNotFoundError(ctx, friendshipconstants.ErrFriendshipNotFound.Error())
		return
	}
	if friendship.ReceiverId != profileId {
		response.RespondForbiddenError(ctx, "You aren't allowed to accept this friend request")
		return
	}

	updateFriendshipErr := friendshipservice.MustUpdateFriendshipStatus(requesterId, profileId, friendshipconstants.StatusAccepted)

	if updateFriendshipErr != nil {
		if errors.Is(updateFriendshipErr, friendshipconstants.ErrFriendshipNotFound) {
			response.RespondNotFoundError(ctx, updateFriendshipErr.Error())
		} else {
			response.RespondInternalError(ctx, updateFriendshipErr.Error())
		}
		return
	}

	response.RespondMessage(ctx, "Friend request has been accepted successfully")
}

func DeleteFriendRequestHandler(ctx *gin.Context) {
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

	friendship, friendshipErr := friendshipservice.MustGetFriendship(profileId, receiverId)

	if friendshipErr != nil {
		if errors.Is(friendshipErr, friendshipconstants.ErrFriendshipNotFound) {
			response.RespondNotFoundError(ctx, friendshipErr.Error())
		} else {
			response.RespondInternalError(ctx, friendshipErr.Error())
		}
		return
	}
	if friendship.RequesterId != profileId && friendship.ReceiverId != profileId {
		response.RespondForbiddenError(ctx, "You aren't allowed to delete this friend request")
		return
	}

	_, deleteFriendshipErr := db.Pool().Exec(context.Background(), "DELETE FROM friendships WHERE (requesterId = $1 AND receiverId = $2) OR (requesterId = $2 AND receiverId = $1);", profileId, receiverId)

	if deleteFriendshipErr != nil {
		response.RespondInternalError(ctx, "Failed to cancel or delete friend request")
		return
	}

	response.RespondMessage(ctx, "Friend request canceled or deleted successfully")
}
