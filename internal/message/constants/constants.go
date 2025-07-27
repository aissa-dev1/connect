package messageconstants

import (
	messagemodel "connect/internal/message/model"
	"errors"
)

var (	
	ErrMessageNotFound = errors.New("message not found")

	ReadStatusSent messagemodel.ReadStatus = "sent"
	ReadStatusDelivered messagemodel.ReadStatus = "delivered"
	ReadStatusSeen messagemodel.ReadStatus = "seen"
	ReadStatusFailed messagemodel.ReadStatus = "failed"
)
