package wshandler

import (
	messageconstants "connect/internal/message/constants"
	messagemodel "connect/internal/message/model"
	messageservice "connect/internal/message/service"
	"connect/internal/pkg/response"
	profileservice "connect/internal/profile/service"
	"connect/internal/ws"
	wsconstants "connect/internal/ws/constants"
	wsmodel "connect/internal/ws/model"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

func HandleChat(ctx *gin.Context) {
	conn, connErr := ws.Upgrader().Upgrade(ctx.Writer, ctx.Request, nil)

	if connErr != nil {
		response.RespondInternalError(ctx, "Failed to upgrade from HTTP to WS")
		ctx.Abort()
		return
	}

	profileId, profileIdErr := profileservice.GetProfileId(ctx)

	if profileIdErr != nil {
		return
	}

	ws.ConnectionManager().Add(profileId, conn)

	defer func() {
		conn.Close()
		ws.ConnectionManager().Remove(profileId)
	}()

	for {
		msgType, msg, readErr := conn.ReadMessage()

		if readErr != nil {
			break
		}

		var data wsmodel.Envelope
		
		if unmarshalErr := json.Unmarshal(msg, &data); unmarshalErr != nil {
			break
		}

		switch data.Type {
			case wsconstants.ChatMessageType:
				var payload wsmodel.ChatMessagePayload

				if unmarshalErr := json.Unmarshal(data.Payload, &payload); unmarshalErr != nil {
					return
				}

				message := messagemodel.Message{
					SenderId: payload.From,
					ReceiverId: payload.To,
					Text: payload.Text,
					ReadStatus: &messageconstants.ReadStatusSent,
				}

				if message.Text == "" {
					break	
				}

				messageId, insertMessageErr := messageservice.InsertMessageAndGetId(message)

				if insertMessageErr != nil {
					return
				}

				message.Id = messageId
				messageJSON, messageMarshalErr := json.Marshal(message)

				if messageMarshalErr != nil {
					return
				}

				envelope := wsmodel.Envelope{
					Type: wsconstants.ChatMessageType,
					Payload: json.RawMessage(messageJSON),
				}

				envelopeJSON, envelopeMarshalErr := json.Marshal(envelope)

				if envelopeMarshalErr != nil {
					return
				}
				if writeErr := conn.WriteMessage(msgType, envelopeJSON); writeErr != nil {
					return
				}

				recipientConn, recipientConnOk := ws.ConnectionManager().Get(payload.To)

				if recipientConnOk {
					if recipientWriteErr := recipientConn.WriteMessage(msgType, envelopeJSON); recipientWriteErr != nil {
						break
					}
				}

			case wsconstants.ChatTypingType:
				var payload wsmodel.ChatTypingPayload

				if unmarshalErr := json.Unmarshal(data.Payload, &payload); unmarshalErr != nil {
					return
				}

				recipientConn, recipientConnOk := ws.ConnectionManager().Get(payload.To)

				if recipientConnOk {
					if recipientWriteErr := recipientConn.WriteMessage(msgType, msg); recipientWriteErr != nil {
						break
					}
				}
		}
	}
}
