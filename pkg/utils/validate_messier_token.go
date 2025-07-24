package utils

import (
	"context"
	"fmt"
	messierTokenRepo "neptune/backend/repositories/messier_token"
	"time"
)

func GetAndValidateMessierToken(ctx context.Context, userID string, repository messierTokenRepo.MessierTokenRepository) (token string, err error) {
	authToken, err := repository.GetMessierTokenByUserID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to get auth token for user %s: %w", userID, err)
	}
	if authToken == nil || authToken.MessierAccessToken == "" {
		return "", fmt.Errorf("no valid auth token found for user %s", userID)
	}

	if authToken.MessierTokenExpires.Before(time.Now()) {
		return "", fmt.Errorf("auth token for user %s has expired", userID)
	}

	return authToken.MessierAccessToken, nil
}
