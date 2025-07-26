package messageconstants

import messagemodel "connect/internal/message/model"

const (
	ReadStatusSent messagemodel.ReadStatus = "sent"
	ReadStatusDelivered messagemodel.ReadStatus = "delivered"
	ReadStatusSeen messagemodel.ReadStatus = "seen"
	ReadStatusFailed messagemodel.ReadStatus = "failed"
)
