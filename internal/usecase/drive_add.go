package usecase

import (
	"context"
	"fmt"
	"github.com/eldius/onedrive-client/client"
	"github.com/eldius/onedrive-client/internal/model"
	"github.com/eldius/onedrive-client/internal/persistence"
	"os"
)

type DriveAddUseUseCase struct {
	c client.Client
	r *persistence.AuthRepository
}

func newDriveAddUseCase(c client.Client, r *persistence.AuthRepository) *DriveAddUseUseCase {
	return &DriveAddUseUseCase{
		c: c,
		r: r,
	}
}

// DriveAdd adds a new drive configuration
func (u *DriveAddUseUseCase) DriveAdd(ctx context.Context, name string) error {
	auth, err := u.c.Authenticate(ctx)
	if err != nil {
		return fmt.Errorf("DriveAdd: authenticate: %w", err)
	}

	fmt.Printf("Auth Token: %s\n", auth.AccessToken)

	user, err := u.c.AuthenticatedUser(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("user: %+v\n", user)

	appDrive, err := u.c.GetAppDriveInfo(ctx)
	if err != nil {
		return fmt.Errorf("DriveAdd: get app drive info: %w", err)
	}

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
	rootFolder, err := u.c.CreateFolder(ctx, hostname, account.Drive.ItemID, account.Drive.DriveID)
	if err != nil {
		return fmt.Errorf("DriveAdd: create folder: %w", err)
	}

	account.Drive.RootFolder = rootFolder.Name

	return u.r.Persist(ctx, account)
}
