package client

import "time"

type APIResponse interface {
	SetStatusCode(int)
	SetRawBody(string)
}

type apiResponse struct {
	RawBody    string
	StatusCode int
}

func (r *apiResponse) SetStatusCode(code int) {
	r.StatusCode = code
}

func (r *apiResponse) SetRawBody(body string) {
	r.RawBody = body
}

type TokenData struct {
	apiResponse
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
}

type CurrentUser struct {
	apiResponse
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

type AppFolderInfo struct {
	apiResponse
	OdataContext         string          `json:"@odata.context"`
	CreatedDateTime      time.Time       `json:"createdDateTime"`
	ETag                 string          `json:"eTag"`
	ID                   string          `json:"id"`
	LastModifiedDateTime time.Time       `json:"lastModifiedDateTime"`
	Name                 string          `json:"name"`
	WebURL               string          `json:"webUrl"`
	CTag                 string          `json:"cTag"`
	Size                 int             `json:"size"`
	CreatedBy            CreatedBy       `json:"createdBy"`
	LastModifiedBy       LastModifiedBy  `json:"lastModifiedBy"`
	ParentReference      ParentReference `json:"parentReference"`
	FileSystemInfo       FileSystemInfo  `json:"fileSystemInfo"`
	Folder               Folder          `json:"folder"`
	SpecialFolder        SpecialFolder   `json:"specialFolder"`
}
type Application struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
}
type User struct {
	Email       string `json:"email"`
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
}
type CreatedBy struct {
	Application Application `json:"application"`
	User        User        `json:"user"`
}
type LastModifiedBy struct {
	Application Application `json:"application"`
	User        User        `json:"user"`
}
type ParentReference struct {
	DriveType string `json:"driveType"`
	DriveID   string `json:"driveId"`
	ID        string `json:"id"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	SiteID    string `json:"siteId"`
}
type FileSystemInfo struct {
	CreatedDateTime      time.Time `json:"createdDateTime"`
	LastModifiedDateTime time.Time `json:"lastModifiedDateTime"`
}
type View struct {
	SortBy    string `json:"sortBy"`
	SortOrder string `json:"sortOrder"`
	ViewType  string `json:"viewType"`
}
type Folder struct {
	ChildCount int  `json:"childCount"`
	View       View `json:"view"`
}
type SpecialFolder struct {
	Name string `json:"name"`
}
