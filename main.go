package main

import (
	"bytes"
	"context"
	"dlfolder/core"
	"dlfolder/sasobjs"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/pkg/browser"
)

// Folder child member query
var memberQuery url.Values = url.Values{}

// Sample query that selects only files who's names end with .sas - modify as needed
func init() {
	memberQuery.Add("filter", "or(and(eq(contentType, 'file'),endsWith(name, '.sas')),eq(contentType, 'folder'))")
}

func main() {
	var baseURL, clientid, clientsecret, rootSourcePath, authcode, rootOutputPath string
	flag.StringVar(&baseURL, "url", "", "Please enter the baseURL")
	flag.StringVar(&clientid, "clientid", "", "Please enter a ClientID")
	flag.StringVar(&rootSourcePath, "folder", "", "Please enter a SAS Content folder path")
	flag.StringVar(&rootOutputPath, "dir", "c:/temp", "Please enter the output directory path")
	flag.Parse()

	if rootSourcePath == "" {
		log.Fatalln("-folder is a required parameter. Aborting.")
	}

	// Open browser to get a one time authorization code
	browser.OpenURL(baseURL + "/SASLogon/oauth/authorize?client_id=" + clientid + "&response_type=code")

	// Get authorization code from end-user assuming he copied it from the browser
	fmt.Println("Enter authorization code displayed in browser:")
	fmt.Scan(&authcode)

	// Get client secret
	fmt.Println("Enter client's secret:")
	fmt.Println("\033[8m") // Hide input
	fmt.Scan(&clientsecret)
	fmt.Println("\033[28m") // Show input

	// Get SAS access token using the auth code
	ai := core.AuthInfo{
		// Username:     username,
		// Password:     password,
		Code:         authcode,
		GrantType:    "authorization_code",
		ClientID:     clientid,
		ClientSecret: clientsecret}

	token := ai.GetToken(baseURL)

	// Store SAS Viya OAuth token and baseURL in a new context
	ctx := context.Background()
	ctx = context.WithValue(ctx, "token", &token)
	ctx = context.WithValue(ctx, "baseURL", baseURL)

	// Folders query see https://developer.sas.com/apis/rest/#making-an-api-call for details on query syntax
	// Search for the specific folder as passed in the input flag
	folderQuery := url.Values{}
	folderQuery.Add("path", rootSourcePath)

	// get root folder ID and call download function
	folderID := sasobjs.GetFolderID(ctx, folderQuery)

	// download
	log.Printf("Downloading %s to directory: %s\n", rootSourcePath, rootOutputPath)
	downloadFolder(ctx, rootOutputPath, folderID)
}

// Handle download of a folder and its child folders/members
func downloadFolder(ctx context.Context, basePath string, folderID string) {

	// get folder details and establish base folder to point to this folder on disk
	folder := sasobjs.GetFolder(ctx, folderID)
	log.Printf("--> Folder Name: %v Id: %v Members: %v\n", folder.Name, folder.ID, folder.MemberCount)

	// update current base path with new folder
	currentBasePath := basePath + "/" + folder.Name + "/"

	// create directory on disk for this folder
	err := os.Mkdir(currentBasePath, 0750)
	if err != nil {
		log.Panicf("Error trying to create output dir %s\n", folder.Name)
	}

	// Get folder members that follow the applied member filter
	// 1. Get the member's (file's) content
	// 2. Save it to a new file in the directory created above
	// 3. For subfolders call this function recursively
	members := sasobjs.GetMembers(ctx, folder.ID, memberQuery)
	for _, member := range members.Items {
		switch member.ContentType {
		case "file":
			downloadMember(ctx, currentBasePath, member.URI, member.Name)
		case "folder":
			downloadFolder(ctx, currentBasePath, strings.Split(member.URI, "/")[3])
		default:
			log.Println("In type other, ignoring.")
		}
	}
}

// Handle download of a folder member
func downloadMember(ctx context.Context, currentBasePath string, memberURI string, memberName string) {
	log.Printf("--> Downloading member Name: %s Member URI: %s\n", memberName, memberURI)
	memberContent := sasobjs.GetFileContent(ctx, memberURI)
	file, err := os.Create(currentBasePath + memberName)
	if err != nil {
		log.Printf("Error trying to create output file %s\n", memberName)
	}
	io.Copy(file, bytes.NewReader(memberContent))
	file.Close()
}
