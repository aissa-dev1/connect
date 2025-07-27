package wsconstants

import wsmodel "connect/internal/ws/model"

const (
	ChatMessageType wsmodel.EnvelopeType = "message"
	ChatTypingType wsmodel.EnvelopeType    = "typing"
)
