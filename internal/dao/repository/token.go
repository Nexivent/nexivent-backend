package repository

import (
	"time"

	"gorm.io/gorm"

	"github.com/Nexivent/nexivent-backend/internal/dao/model"
	"github.com/Nexivent/nexivent-backend/logging"
)

type Token struct {
	logger logging.Logger
	DB     *gorm.DB
}

func (m Token) New(userID int64, ttl time.Duration, scope string) (*model.Token, error) {
	token, err := model.GenerateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}
	err = m.Insert(token)
	return token, err
}

func (m Token) Insert(token *model.Token) error {
	result := m.DB.Create(token)
	if result.Error != nil {
		m.logger.Errorf("Token.Insert: %v", result.Error)
		return result.Error
	}
	return nil
}

func (m Token) DeleteAllForUser(scope string, userID int64) error {
	result := m.DB.Where("usuario_id = ? AND scope = ?", userID, scope).Delete(&model.Token{})
	if result.Error != nil {
		m.logger.Errorf("Token.DeleteAllForUser: %v", result.Error)
		return result.Error
	}
	return nil
}
