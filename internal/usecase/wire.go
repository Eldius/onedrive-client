//go:build wireinject

package usecase

import (
	"github.com/eldius/onedrive-client/client"
	"github.com/eldius/onedrive-client/internal/persistence"
	"github.com/google/wire"
)

func NewFileUpload(_ client.Client) *FileUploadUseCase {
	wire.Build(persistence.NewAuthRepository, persistence.NewDB, newFileUploadUseCase)
	return nil
}

func NewListFilesUseCase(_ client.Client) *ListFilesUseCase {
	wire.Build(persistence.NewAuthRepository, persistence.NewDB, newListFilesUseCase)
	return nil
}

func NewDriveAddUseUseCase(_ client.Client) *DriveAddUseUseCase {
	wire.Build(persistence.NewAuthRepository, persistence.NewDB, newDriveAddUseCase)
	return nil
}
