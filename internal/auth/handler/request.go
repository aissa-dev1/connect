package authhandler

import (
	authmodel "connect/internal/auth/model"
	authservice "connect/internal/auth/service"
	"connect/internal/db"
	errormessage "connect/internal/pkg/error_message"
	"connect/internal/pkg/hasher"
	"connect/internal/pkg/response"
	usermodel "connect/internal/user/model"
	userservice "connect/internal/user/service"
	"connect/utils"
	"context"

	"github.com/gin-gonic/gin"
)

func SignUpHandler(ctx *gin.Context) {
	var requestBody authmodel.SignUpRequestBody

	requestBodyErr := ctx.BindJSON(&requestBody)

	if requestBodyErr != nil {
		response.RespondInternalError(ctx, errormessage.InternalServerError)
		return
	}
	if !utils.ValidEmail(requestBody.Email) {
		response.RespondBadRequestError(ctx, errormessage.EnterValidEmail)
		return
	}
	if requestBody.Username == "" {
		response.RespondBadRequestError(ctx, "Please provide a username")
		return
	}
	if len(requestBody.Password) < 4 {
		response.RespondBadRequestError(ctx, errormessage.MustHaveStrongPassword)
		return
	}

	userExists, userExistsErr := userservice.UserExists(requestBody.Email)

	if userExistsErr != nil {
		response.RespondInternalError(ctx, userExistsErr.Error())
		return
	}
	if userExists {
		response.RespondConflictError(ctx, "This email is already linked with another account")
		return
	}

	hashedPassword, hashedPasswordErr := hasher.Global().Hash(requestBody.Password)

	if hashedPasswordErr != nil {
		response.RespondInternalError(ctx, hashedPasswordErr.Error())
		return
	}

	insertUserErr := userservice.InsertUser(usermodel.User{Email: requestBody.Email, Username: requestBody.Username, Password: hashedPassword})

	if insertUserErr != nil {
		response.RespondInternalError(ctx, errormessage.InternalServerError)
		return
	}

	var userId int

	userIdErr := db.Pool().QueryRow(context.Background(), "SELECT id FROM users WHERE email = $1;", requestBody.Email).Scan(&userId)

	if userIdErr != nil {
		response.RespondInternalError(ctx, errormessage.InternalServerError)
		return
	}

	_, tokenErr := authservice.RegisterToken(ctx, userId)

	if tokenErr != nil {
		response.RespondInternalError(ctx, tokenErr.Error())
		return
	}

	response.RespondMessage(ctx, "Welcome "+requestBody.Username)
}

func SignInHandler(ctx *gin.Context) {
	var requestBody authmodel.SignInRequestBody

	requestBodyErr := ctx.BindJSON(&requestBody)

	if requestBodyErr != nil {
		response.RespondInternalError(ctx, errormessage.InternalServerError)
		return
	}
	if !utils.ValidEmail(requestBody.Email) {
		response.RespondBadRequestError(ctx, errormessage.EnterValidEmail)
		return
	}
	if len(requestBody.Password) < 4 {
		response.RespondBadRequestError(ctx, errormessage.MustHaveStrongPassword)
		return
	}

	var user usermodel.MinimalUserWithPassword

	userErr := db.Pool().QueryRow(context.Background(), "SELECT id, username, password FROM users WHERE email = $1;", requestBody.Email).Scan(&user.Id, &user.Username, &user.Password)

	if userErr != nil {
		response.RespondNotFoundError(ctx, errormessage.InvalidCredentials)
		return
	}

	passwordMatch := hasher.Global().Compare(requestBody.Password, user.Password)

	if !passwordMatch {
		response.RespondNotFoundError(ctx, errormessage.InvalidCredentials)
		return
	}

	_, tokenErr := authservice.RegisterToken(ctx, user.Id)

	if tokenErr != nil {
		response.RespondInternalError(ctx, tokenErr.Error())
		return
	}

	response.RespondMessage(ctx, "Welcome back "+user.Username)
}

func SignOutHandler(ctx *gin.Context) {
	authservice.ClearToken(ctx)
	response.RespondMessage(ctx, "Signed out successfully")
}
