package types

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

type CreateFile struct {
	apiResponse
	OdataContext              string          `json:"@odata.context,omitempty"`
	OdataEtag                 string          `json:"@odata.etag,omitempty"`
	MicrosoftGraphDownloadURL string          `json:"@microsoft.graph.downloadUrl,omitempty"`
	CreatedDateTime           time.Time       `json:"createdDateTime,omitempty"`
	ETag                      string          `json:"eTag,omitempty"`
	ID                        string          `json:"id,omitempty"`
	LastModifiedDateTime      time.Time       `json:"lastModifiedDateTime,omitempty"`
	Name                      string          `json:"name,omitempty"`
	Size                      int             `json:"size,omitempty"`
	WebURL                    string          `json:"webUrl,omitempty"`
	CTag                      string          `json:"cTag,omitempty"`
	CommentSettings           CommentSettings `json:"commentSettings,omitempty"`
	CreatedBy                 CreatedBy       `json:"createdBy,omitempty"`
	LastModifiedBy            LastModifiedBy  `json:"lastModifiedBy,omitempty"`
	ParentReference           ParentReference `json:"parentReference,omitempty"`
	File                      File            `json:"file,omitempty"`
	FileSystemInfo            FileSystemInfo  `json:"fileSystemInfo,omitempty"`
	Shared                    Shared          `json:"shared,omitempty"`
}
type CommentingDisabled struct {
	IsDisabled bool `json:"isDisabled,omitempty"`
}
type CommentSettings struct {
	CommentingDisabled CommentingDisabled `json:"commentingDisabled,omitempty"`
}
type Application struct {
	DisplayName string `json:"displayName,omitempty"`
	ID          string `json:"id,omitempty"`
}
type User struct {
	DisplayName string `json:"displayName,omitempty"`
	Email       string `json:"email,omitempty"`
}
type CreatedBy struct {
	Application Application `json:"application,omitempty"`
	User        User        `json:"user,omitempty"`
}
type LastModifiedBy struct {
	Application Application `json:"application,omitempty"`
	User        User        `json:"user,omitempty"`
}
type SharepointIds struct {
	ListID           string `json:"listId,omitempty"`
	ListItemUniqueID string `json:"listItemUniqueId,omitempty"`
	SiteID           string `json:"siteId,omitempty"`
	SiteURL          string `json:"siteUrl,omitempty"`
	TenantID         string `json:"tenantId,omitempty"`
	WebID            string `json:"webId,omitempty"`
}
type ParentReference struct {
	DriveID       string        `json:"driveId,omitempty"`
	DriveType     string        `json:"driveType,omitempty"`
	ID            string        `json:"id,omitempty"`
	Path          string        `json:"path,omitempty"`
	SharepointIds SharepointIds `json:"sharepointIds,omitempty"`
}
type Hashes struct {
	QuickXorHash string `json:"quickXorHash,omitempty"`
}
type File struct {
	MimeType string `json:"mimeType,omitempty"`
	Hashes   Hashes `json:"hashes,omitempty"`
}
type FileSystemInfo struct {
	CreatedDateTime      time.Time `json:"createdDateTime,omitempty"`
	LastModifiedDateTime time.Time `json:"lastModifiedDateTime,omitempty"`
}
type Shared struct {
	Scope string `json:"scope,omitempty"`
}

type ListFiles struct {
	apiResponse
	OdataContext string  `json:"@odata.context"`
	Value        []Value `json:"value"`
}
type Owner struct {
	User User `json:"user"`
}
type Value struct {
	MicrosoftGraphDownloadURL string          `json:"@microsoft.graph.downloadUrl"`
	CreatedDateTime           time.Time       `json:"createdDateTime"`
	ETag                      string          `json:"eTag"`
	ID                        string          `json:"id"`
	LastModifiedDateTime      time.Time       `json:"lastModifiedDateTime"`
	Name                      string          `json:"name"`
	WebURL                    string          `json:"webUrl"`
	CTag                      string          `json:"cTag"`
	Size                      int             `json:"size"`
	CreatedBy                 CreatedBy       `json:"createdBy"`
	LastModifiedBy            LastModifiedBy  `json:"lastModifiedBy"`
	ParentReference           ParentReference `json:"parentReference"`
	File                      File            `json:"file"`
	FileSystemInfo            FileSystemInfo  `json:"fileSystemInfo"`
	Shared                    Shared          `json:"shared"`
}
