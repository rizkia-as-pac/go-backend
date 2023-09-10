package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/tech_school/simple_bank/db/sqlc"
)

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	// AccessTokenExpiresAt akan berguna bagi client untuk setup shcedule to renew accesstoken
	AccessTokenExpireAt time.Time `json:"access_token_expired_at"`
}

func (server *Server) renewAccessToken(ctx *gin.Context) {
	var req renewAccessTokenRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	RTPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	session, err := server.store.GetSession(ctx, RTPayload.ID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if session.IsBlocked {
		err := fmt.Errorf("blocked session")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.Username != RTPayload.Username {
		err := fmt.Errorf("incorrect session user")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.RefreshToken != req.RefreshToken {
		err := fmt.Errorf("mismatched session token")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if time.Now().After(session.ExpiredAt) {
		err := fmt.Errorf("expired session")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, ATPayload, err := server.tokenMaker.CreateToken(RTPayload.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := renewAccessTokenResponse{
		AccessToken:         accessToken,
		AccessTokenExpireAt: ATPayload.ExpiredAt,
	}

	ctx.JSON(http.StatusOK, res)

}
