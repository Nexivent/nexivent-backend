package controller

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/dao/repository"
	"github.com/Nexivent/nexivent-backend/logging"
	"github.com/Nexivent/nexivent-backend/internal/dao/model"
	"github.com/Nexivent/nexivent-backend/errors"
)

type TokenController struct {
	Logger logging.Logger
	DB     *repository.NexiventPsqlEntidades
}

func (tc *TokenController) CreateToken(userID int64, ttl time.Duration, scope string) (*model.Token, *errors.Error) {
	token, err := tc.DB.Token.New(userID, ttl, scope)
	if err != nil {
		return nil, &errors.InternalServerError.Default
	}
	return token, nil
}

func (tc *TokenController) ValidateToken(tokenValue string) (*model.Token, *errors.Error) {
	var token model.Token
	result := tc.DB.Token.DB.Where("value = ?", tokenValue).First(&token)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, &errors.AuthenticationError.InvalidAccessToken
		}
		tc.Logger.Errorf("TokenController.ValidateToken: %v", result.Error)
		return nil, &errors.InternalServerError.Default
	}

	if time.Now().After(token.Expiry) {
		return nil, &errors.AuthenticationError.ExpiredToken
	}

	return &token, nil
}

// func (tc *TokenController) DeleteTokensForUser(scope string, userID int64) *logging.Error {
// 	err := tc.DB.Token.DeleteAllForUser(scope, userID)
// 	if err != nil {
// 		tc.Logger.Errorf("TokenController.DeleteTokensForUser: %v", err)
// 		return logging.NewErrorFromError(err, "TOKEN_ERROR_002", "Failed to delete tokens for user")
// 	}
// 	return nil
// }