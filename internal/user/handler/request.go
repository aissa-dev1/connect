package userhandler

import (
	"connect/internal/db"
	errormessage "connect/internal/pkg/error_message"
	"connect/internal/pkg/response"
	usermodel "connect/internal/user/model"
	"context"

	"github.com/gin-gonic/gin"
)

func SearchUsersHandler(ctx *gin.Context) {
	users := []usermodel.SearchUser{}

	userNameQuery := ctx.Query("q")

	if len(userNameQuery) > 0 {
		userNameQuery = "%" + userNameQuery + "%"
	}

	rows, rowsErr := db.Pool().Query(context.Background(), "SELECT username FROM users WHERE username ILIKE $1;", userNameQuery)

	if rowsErr != nil {
		response.RespondInternalError(ctx, errormessage.InternalServerError)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var user usermodel.SearchUser

		scanErr := rows.Scan(&user.Username)

		if scanErr != nil {
			response.RespondInternalError(ctx, errormessage.InternalServerError)
			return
		}

		users = append(users, user)
	}

	response.RespondData[[]usermodel.SearchUser](ctx, users)
}
