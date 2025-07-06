package messier_token

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	model "neptune/backend/models/user"
)

type messierTokenRepositoryImplement struct {
	db *gorm.DB
}

func (m messierTokenRepositoryImplement) Save(ctx context.Context, token *model.MessierToken) error {
	// GORM's Upsert for PostgreSQL/SQLite
	err := m.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}}, // Conflict on UserID
		DoUpdates: clause.Assignments(map[string]interface{}{ // Update these fields on conflict
			"messier_access_token":  token.MessierAccessToken,
			"messier_token_expires": token.MessierTokenExpires,
			// GORM handles `updated_at` automatically with `gorm.Model` on updates
		}),
	}).Create(token).Error // Use Create method with Clauses
	if err != nil {
		return fmt.Errorf("failed to save or update external token for user %s: %w", token.UserID, err)
	}
	return nil
}

func (m messierTokenRepositoryImplement) GetMessierTokenByUserID(ctx context.Context, userID string) (*model.MessierToken, error) {
	var token model.MessierToken
	result := m.db.WithContext(ctx).Where("user_id = ?", userID).First(&token)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Not found is not an error, just return nil
		}
		return nil, fmt.Errorf("failed to find external token by user ID %s: %w", userID, result.Error)
	}
	return &token, nil
}

func (m messierTokenRepositoryImplement) DeleteByUserID(ctx context.Context, userID string) error {
	result := m.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&model.MessierToken{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete external token for user ID %s: %w", userID, result.Error)
	}
	return nil
}

func NewMessierTokenRepository(db *gorm.DB) MessierTokenRepository {
	return &messierTokenRepositoryImplement{db: db}
}
