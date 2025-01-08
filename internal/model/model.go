package model

type OnedriveAccount struct {
	ID       string `gorm:"id"`
	Name     string `gorm:"index"`
	AuthData *TokenData
}

type TokenData struct {
	TokenType    string
	Scope        string
	ExpiresIn    int
	ExtExpiresIn int
	AccessToken  string
	RefreshToken string
	IDToken      string
}
