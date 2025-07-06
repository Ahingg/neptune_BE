package messier_token

import (
	"context"
	model "neptune/backend/models/user"
)

type MessierTokenRepository interface {
	Save(ctx context.Context, token *model.MessierToken) error
	GetMessierTokenByUserID(ctx context.Context, userID string) (*model.MessierToken, error)
	DeleteByUserID(ctx context.Context, userID string) error
}
