package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"taxApi/internal/models"
	"time"
)

func (s *Service) CreatingToken(token string) (string, error) {
	buffer := make([]byte, 256)
	n, err := rand.Read(buffer)
	if n != len(buffer) {
		return "", errors.New("не удалось сгенерировать байт")
	}
	if err != nil {
		return "", err
	}
	token = hex.EncodeToString(buffer)
	return token, nil
}

func (s *Service) GetTokenToUser(ctx context.Context, user *models.User) (token string, err error) {
	user, err = s.ValidateLoginAndPass(user.Login, user.Password)
	if err != nil {
		return "", err
	}

	token, err = s.CreatingToken(token)
	if err != nil {
		return "", err
	}

	err = s.Repository.PutNewToken(ctx, token, int64(user.Id))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) IdByToken(ctx context.Context, token string) (id int, err error) {
	id, expire, err := s.Repository.IdByToken(ctx, token)
	if err != nil {
		return 0, err
	}
	
	if time.Now().After(expire) {
		return 0, errors.New("expired")
	}
	return id, nil

}
