package usecase

import (
	"context"
	"fmt"
	"github.com/eldius/onedrive-client/internal/model"
	"github.com/eldius/onedrive-client/internal/persistence"
	"log/slog"
)

func loadSession(ctx context.Context, r *persistence.AuthRepository, accName string) (*model.OnedriveAccount, error) {
	acc, err := r.FindOneByName(ctx, accName)
	if err != nil {
		return nil, fmt.Errorf("could not find account %q: %w", accName, err)
	}
	slog.With("account_name", accName, "account", acc).Info("found account")
	return acc, err
}
