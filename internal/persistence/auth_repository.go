package persistence

import (
	"context"
	"fmt"
	"github.com/eldius/onedrive-client/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) Persist(ctx context.Context, a *model.OnedriveAccount) error {
	if a.ID == "" {
		a.ID = uuid.NewString()
	}
	tx := r.db.WithContext(ctx).Save(a)
	if tx.Error != nil {
		return fmt.Errorf("save onedrive account: %w", tx.Error)
	}
	return nil
}

func (r *AuthRepository) FindOneByName(ctx context.Context, name string) (*model.OnedriveAccount, error) {
	var acc model.OnedriveAccount
	if tx := r.db.WithContext(ctx).First(&acc, "Name", name); tx.Error != nil {
		return nil, fmt.Errorf("find onedrive account: %w", tx.Error)
	}
	var drv model.DriveInfo
	if tx := r.db.WithContext(ctx).First(&drv, "account_id", acc.ID); tx.Error != nil {
		return nil, fmt.Errorf("find onedrive account drive info: %w", tx.Error)
	}
	acc.Drive = &drv

	var auth model.TokenData
	if tx := r.db.WithContext(ctx).First(&auth, "account_id", acc.ID); tx.Error != nil {
		return nil, fmt.Errorf("find onedrive account auth info: %w", tx.Error)
	}
	acc.AuthData = &auth

	return &acc, nil
}
