package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/tech_school/simple_bank/db/sqlc"
	"github.com/tech_school/simple_bank/utils/pass"
)

type createUserRequest struct {
	// alphanum = membuat user tidak dapat menginputkan karakter untuk username selain alphabet dan number
	Username string `json:"username" binding:"required,alphanum"`
	// panjang password sebaiknya tidak terlalu pendek
	Password string `json:"password" binding:"required,min=4"`
	FullName string `json:"full_name" binding:"required"`
	// email = menvalidasi agar data yang diinputkan user berformat email yang benar
	Email string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Username          string    `json:"username"`
	Fullname          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {

	return userResponse{
		Username:          user.Username,
		Fullname:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := pass.HashedPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := newUserResponse(user)

	ctx.JSON(http.StatusOK, res)

}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=4"`
}

type loginUserResponse struct {
	SessionID   uuid.UUID `json:"session_id"`
	AccessToken string    `json:"access_token"`
	// AccessTokenExpiresAt akan berguna bagi client untuk setup shcedule to renew accesstoken
	AccessTokenExpireAt  time.Time `json:"access_token_expired_at"`
	RefreshToken         string    `json:"refresh_token"`
	RefreshTokenExpireAt time.Time `json:"refresh_token_expired_at"`
	User                 userResponse
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = pass.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, ATPayload, err := server.tokenMaker.CreateToken(req.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, RTPayload, err := server.tokenMaker.CreateToken(req.Username, server.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateSessionParams{
		ID:           RTPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiredAt:    RTPayload.ExpiredAt,
	}
	session, err := server.store.CreateSession(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := loginUserResponse{
		SessionID:            session.ID,
		AccessToken:          accessToken,
		AccessTokenExpireAt:  ATPayload.ExpiredAt,
		RefreshToken:         refreshToken,
		RefreshTokenExpireAt: RTPayload.ExpiredAt,
		User:                 newUserResponse(user),
	}

	ctx.JSON(http.StatusOK, res)

}
