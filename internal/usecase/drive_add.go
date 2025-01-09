package usecase

import (
	"context"
	"fmt"
	"github.com/eldius/onedrive-client/client"
	"github.com/eldius/onedrive-client/internal/configs"
	"github.com/eldius/onedrive-client/internal/model"
	"github.com/eldius/onedrive-client/internal/persistence"
	"os"
)

// DriveAdd adds a new drive configuration
func DriveAdd(ctx context.Context, name string) error {
	c := client.New(
		client.WithSecretID(configs.GetSecretID()),
		client.WithScopes(configs.GetAuthScopes()...),
		client.WithRedirectURL(configs.GetRedirectURL()),
	)

	auth, err := c.Authenticate(ctx)
	if err != nil {
		return fmt.Errorf("DriveAdd: authenticate: %w", err)
	}

	fmt.Printf("Auth Token: %s\n", auth.AccessToken)

	user, err := c.AuthenticatedUser(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("user: %+v\n", user)

	appDrive, err := c.GetAppDriveInfo(ctx)
	if err != nil {
		return fmt.Errorf("DriveAdd: get app drive info: %w", err)
	}

	r := persistence.NewAuthRepository(persistence.GetDB())

	account := &model.OnedriveAccount{
		Name: name,
		AuthData: &model.TokenData{
			TokenType:    auth.TokenType,
			Scope:        auth.Scope,
			ExpiresIn:    auth.ExpiresIn,
			ExtExpiresIn: auth.ExtExpiresIn,
			AccessToken:  auth.AccessToken,
			RefreshToken: auth.RefreshToken,
			IDToken:      auth.IDToken,
		},
		Drive: &model.DriveInfo{
			DriveID: appDrive.ParentReference.DriveID,
			ItemID:  appDrive.ID,
		},
	}

	hostname, _ := os.Hostname()
	rootFolder, err := c.CreateFolder(ctx, hostname, account.Drive.ItemID, account.Drive.DriveID)
	if err != nil {
		return fmt.Errorf("DriveAdd: create folder: %w", err)
	}

	account.Drive.RootFolder = rootFolder.Name

	return r.Persist(ctx, account)
}
