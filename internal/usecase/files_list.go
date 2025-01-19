package usecase

import (
	"context"
	"fmt"
	"github.com/eldius/onedrive-client/client"
	"github.com/eldius/onedrive-client/client/types"
	"github.com/eldius/onedrive-client/internal/persistence"
)

type ListFilesUseCase struct {
	c client.Client
	r *persistence.AuthRepository
}

func newListFilesUseCase(r *persistence.AuthRepository) *ListFilesUseCase {
	return &ListFilesUseCase{
		r: r,
	}
}

func (l *ListFilesUseCase) ListFilesFromDrive(ctx context.Context, accountName string) error {
	acc, err := loadSession(ctx, l.r, accountName)
	if err != nil {
		return fmt.Errorf("could not find account %q: %w", accountName, err)
	}

	token := &types.TokenData{
		TokenType:    acc.AuthData.TokenType,
		Scope:        acc.AuthData.Scope,
		ExpiresIn:    acc.AuthData.ExpiresIn,
		ExtExpiresIn: acc.AuthData.ExtExpiresIn,
		AccessToken:  acc.AuthData.AccessToken,
		RefreshToken: acc.AuthData.RefreshToken,
		IDToken:      acc.AuthData.IDToken,
	}
	c := client.New(
		client.WithScopes(acc.AuthData.Scope),
		client.WithAuthenticationTokenData(token),
	)
	remoteFiles, err := c.ListFiles(ctx, acc.Drive.DriveID, acc.Drive.ItemID)
	if err != nil {
		return fmt.Errorf("listing files: %w", err)
	}

	for _, f := range remoteFiles.Value {
		fmt.Printf(" -> file: %s (%s)\n", f.Name, f.GetMimeType())
	}

	return nil
}
