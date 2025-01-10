package model

import "time"

type OnedriveAccount struct {
	ID        string     `gorm:"id"`
	Name      string     `gorm:"index"`
	AuthData  *TokenData `gorm:"foreignKey:AccountID"`
	Drive     *DriveInfo `gorm:"foreignKey:AccountID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TokenData struct {
	TokenType    string
	Scope        string
	ExpiresIn    int
	ExtExpiresIn int
	AccessToken  string
	RefreshToken string
	IDToken      string
	AccountID    string `gorm:"index"`
}

type DriveInfo struct {
	ID         string `gorm:"id"`
	DriveID    string `gorm:"index"`
	ItemID     string `gorm:"index"`
	RootFolder string `gorm:"index"`
	AccountID  string `gorm:"index"`
}
