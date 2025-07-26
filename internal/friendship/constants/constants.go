package friendshipconstants

import (
	friendshipmodel "connect/internal/friendship/model"
	"errors"
)

var (
	ErrFriendshipNotFound = errors.New("friendship not found")

	StatusPending friendshipmodel.FriendshipStatus = "pending"
	StatusAccepted friendshipmodel.FriendshipStatus = "accepted"
	StatusRejected friendshipmodel.FriendshipStatus = "rejected"
)
