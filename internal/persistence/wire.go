//go:build wireinject

package persistence

import (
	"github.com/eldius/onedrive-client/internal/configs"
	"github.com/google/wire"
	"gorm.io/gorm"
)

func NewDB() *gorm.DB {
	wire.Build(configs.GetDBFilePath, GetDB)
	return nil
}
