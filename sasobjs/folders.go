package sasobjs

import (
	"context"
	"dlfolder/core"
	"encoding/json"
	"net/url"
)

// FolderList is an object representing a list of folders
type FolderList struct {
	Name    string   `json:"name"`
	Start   int      `json:"start"`
	Limit   int      `json:"limit"`
	Count   int      `json:"count"`
	Accept  string   `json:"accept"`
	Links   []Link   `json:"links"`
	Version int      `json:"version"`
	Items   []Folder `json:"items"`
}

// Folder is a folder object
type Folder struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	ParentFolderURI   string `json:"parentFolderUri"`
	CreationTimeStamp string `json:"creationTimeStamp"`
	ModifiedTimeStamp string `json:"modifiedTimeStamp"`
	CreatedBy         string `json:"createdBy"`
	ModifiedBy        string `json:"modifiedBy"`
	Type              string `json:"type"`
	IconURI           string `json:"iconUri"`
	MemberCount       int    `json:"memberCount"`
	Links             []Link `json:"links"`
	Properties        string `json:"properties"`
	Version           int    `json:"version"`
}

// MemberList is an object representing a list of members in a folder
type MemberList struct {
	Name    string   `json:"name"`
	Start   int      `json:"start"`
	Limit   int      `json:"limit"`
	Count   int      `json:"count"`
	Accept  string   `json:"accept"`
	Links   []Link   `json:"links"`
	Version int      `json:"version"`
	Items   []Member `json:"items"`
}

// Member is a member object
type Member struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	ParentFolderURI   string `json:"parentFolderUri"`
	CreationTimeStamp string `json:"creationTimeStamp"`
	ModifiedTimeStamp string `json:"modifiedTimeStamp"`
	CreatedBy         string `json:"createdBy"`
	ModifiedBy        string `json:"modifiedBy"`
	Type              string `json:"type"`
	ContentType       string `json:"contentType"`
	URI               string `json:"uri"`
	Links             []Link `json:"links"`
	Version           int    `json:"version"`
}

// Link is a link object
type Link struct {
	Method           string `json:"method"`
	Rel              string `json:"rel"`
	URI              string `json:"uri"`
	Href             string `json:"href"`
	Title            string `json:"title"`
	Type             string `json:"type"`
	ItemType         string `json:"itemType"`
	ResponseType     string `json:"responseType"`
	ResponseItemType string `json:"responseItemType"`
}

// GetFolders extract a list of folders extra filters can be applied
func GetFolders(ctx context.Context, query url.Values) FolderList {
	// bearer := "Bearer " + token.AccessToken
	bearer := "Bearer " + ctx.Value("accessToken").(*core.Token).AccessToken
	baseURL := ctx.Value("baseURL").(string)
	headers := map[string][]string{
		"Accept":        []string{"application/vnd.sas.collection+json"},
		"Authorization": []string{bearer}}
	endpoint := "/folders/folders"
	method := "GET"
	resp := core.CallRest(baseURL, endpoint, headers, method, nil, query)
	var result FolderList
	json.Unmarshal(resp, &result)
	return result
}

// GetFolders extract a list of folders extra filters can be applied
func GetMembers(ctx context.Context, folderid string, query url.Values) MemberList {
	// bearer := "Bearer " + token.AccessToken
	bearer := "Bearer " + ctx.Value("accessToken").(*core.Token).AccessToken
	baseURL := ctx.Value("baseURL").(string)
	headers := map[string][]string{
		"Accept":        []string{"application/vnd.sas.collection+json"},
		"Authorization": []string{bearer}}
	endpoint := "/folders/folders/" + folderid + "/members"
	method := "GET"
	resp := core.CallRest(baseURL, endpoint, headers, method, nil, query)
	var result MemberList
	json.Unmarshal(resp, &result)
	return result
}

// GetFileContent downloads the file as a slice of bytes
func GetFileContent(ctx context.Context, fileurl string, query url.Values) []byte {
	// bearer := "Bearer " + token.AccessToken
	bearer := "Bearer " + ctx.Value("accessToken").(*core.Token).AccessToken
	baseURL := ctx.Value("baseURL").(string)
	headers := map[string][]string{
		// "Accept":        []string{"application/vnd.sas.collection+json"},
		"Authorization": []string{bearer}}
	endpoint := fileurl + "/content"
	method := "GET"
	resp := core.CallRest(baseURL, endpoint, headers, method, nil, query)
	return resp
}