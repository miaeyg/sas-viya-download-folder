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

	"github.com/pkg/browser"
)

func main() {
	var hostname, clientid, clientsecret, search, limit, code, output string
	flag.StringVar(&output, "o", "c:/temp", "Please enter the output path")
	flag.StringVar(&hostname, "h", "", "Please enter the hostname")
	flag.StringVar(&clientid, "ci", "", "Please enter a ClientID")
	flag.StringVar(&search, "s", "", "Please enter a search string")
	flag.StringVar(&limit, "l", "100", "Please enter a limit on the returned number of results")
	flag.Parse()

	// open browser to get authorization code
	browser.OpenURL(hostname + "/SASLogon/oauth/authorize?client_id=" + clientid + "&response_type=code")

	// get authorization code
	fmt.Println("Enter code from browser:")
	fmt.Scan(&code)

	// get client secret
	fmt.Println("Enter client secret:")
	fmt.Println("\033[8m") // Hide input
	fmt.Scan(&clientsecret)
	fmt.Println("\033[28m") // Show input

	ai := core.AuthInfo{
		// Username:     username,
		// Password:     password,
		Code:         code,
		GrantType:    "authorization_code",
		ClientID:     clientid,
		ClientSecret: clientsecret}

	baseURL := hostname
	token := ai.GetToken(baseURL)

	// store token and baseURL in a context
	ctx := context.Background()
	ctx = context.WithValue(ctx, "token", &token)
	ctx = context.WithValue(ctx, "baseURL", baseURL)

	// folders query see https://developer.sas.com/apis/rest/#making-an-api-call for details on query syntax
	queryfl := url.Values{}
	queryfl.Add("limit", limit)
	if search != "" {
		queryfl.Add("filter", "contains(name, "+search+")")
	}

	// members query - only files that contain ".sas"
	querymem := url.Values{}
	querymem.Add("filter", "and(eq(contentType, 'file'),endsWith(name, '.sas'))")

	log.Println("Root output folder:", output)
	log.Println("Searching for folders in SAS Content...")

	// get list of folders and loop over them
	folders := sasobjs.GetFolders(ctx, queryfl)
	for _, folder := range folders.Items {

		log.Printf("Folder Name: %v Id: %v Members: %v\n", folder.Name, folder.ID, folder.MemberCount)
		err := os.Mkdir(output+"/"+folder.Name, 0750)
		if err != nil {
			fmt.Println("Error trying to create output dir", folder.Name)
			continue
		}

		// get folder members that follow the member filter
		members := sasobjs.GetMembers(ctx, folder.ID, querymem)
		for _, member := range members.Items {
			log.Printf("Downloading member Name: %s Member URI: %s Member ID: %s\n", member.Name, member.URI, member.ID)
			sasfile := sasobjs.GetFileContent(ctx, member.URI)
			file, err := os.Create(output + "/" + folder.Name + "/" + member.Name)
			if err != nil {
				fmt.Println("Error trying to create output file", member.Name)
				continue
			}
			io.Copy(file, bytes.NewReader(sasfile))
			file.Close()
		}
	}
}
