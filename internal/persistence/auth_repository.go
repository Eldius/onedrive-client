package persistence

import (
	"context"
	"fmt"
	"github.com/eldius/onedrive-client/internal/model"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) Persist(ctx context.Context, a *model.OnedriveAccount) error {
	tx := r.db.WithContext(ctx).Save(a)
	if tx.Error != nil {
		return fmt.Errorf("save onedrive account: %w", tx.Error)
	}
	return nil
}
