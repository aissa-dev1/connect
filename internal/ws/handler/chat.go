package wshandler

import (
	errormessage "connect/internal/pkg/error_message"
	"connect/internal/pkg/response"
	"connect/internal/ws"
	"fmt"

	"github.com/gin-gonic/gin"
)

func HandleChat(ctx *gin.Context) {
	conn, connErr := ws.Upgrader().Upgrade(ctx.Writer, ctx.Request, nil)

	if connErr != nil {
		response.RespondInternalError(ctx, errormessage.InternalServerError)
		ctx.Abort()
		return
	}

	defer conn.Close()

	for {
		msgType, msg, readErr := conn.ReadMessage()

		if readErr != nil {
			response.RespondInternalError(ctx, errormessage.EnterValidEmail)
			ctx.Abort()
			break
		}

		fmt.Printf("Received: %s\n", msg)

		writeErr := conn.WriteMessage(msgType, msg)

		if writeErr != nil {
			response.RespondInternalError(ctx, errormessage.EnterValidEmail)
			ctx.Abort()
			break
		}
	}
}
