package usecase

import (
	"context"
	"github.com/eldius/onedrive-client/client"
	"github.com/eldius/onedrive-client/internal/persistence"
)

type FileUploadUseCase struct {
	c client.Client
	r *persistence.AuthRepository
}

func newFileUploadUseCase(c client.Client, r *persistence.AuthRepository) *FileUploadUseCase {
	return &FileUploadUseCase{
		c: c,
		r: r,
	}
}

func (u *FileUploadUseCase) Upload(ctx context.Context, accName, inputFile, outputFile string) error {
	//acc, err := loadSession(ctx, u.r, accName)
	//if err != nil {
	//	return fmt.Errorf("loadSession: %w", err)
	//}

	return nil
}
