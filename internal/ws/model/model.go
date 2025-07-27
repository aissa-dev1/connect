package wsmodel

import "encoding/json"

type EnvelopeType = string

type Envelope struct {
	Type    EnvelopeType `json:"type"`
	Payload json.RawMessage   `json:"payload"`
}

type ChatMessagePayload struct {
	From     int  `json:"from"`
	To       int  `json:"to"`
	Text string `json:"text"`
}

type ChatTypingPayload struct {
	From     int  `json:"from"`
	To       int  `json:"to"`
	RecipientUserName string `jsong:"recipientUserName"`
	IsTyping bool `json:"isTyping"`
}
