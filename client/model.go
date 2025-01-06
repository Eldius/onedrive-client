package client

type TokenData struct {
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
}

type CurrentUser struct {
	OdataContext      string   `json:"@odata.context"`
	UserPrincipalName string   `json:"userPrincipalName"`
	ID                string   `json:"id"`
	DisplayName       string   `json:"displayName"`
	Surname           string   `json:"surname"`
	GivenName         string   `json:"givenName"`
	PreferredLanguage string   `json:"preferredLanguage"`
	Mail              string   `json:"mail"`
	MobilePhone       string   `json:"mobilePhone"`
	JobTitle          string   `json:"jobTitle"`
	OfficeLocation    string   `json:"officeLocation"`
	BusinessPhones    []string `json:"businessPhones"`
}
