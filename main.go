package main

import (
	"dlfolder/core"
	"dlfolder/sasobjs"
	"flag"
	"fmt"
	"net/url"
)

func main() {
	var username, password, hostname, clientid, clientsecret, search, limit, code string
	flag.StringVar(&username, "u", "", "Please enter a user name")
	flag.StringVar(&password, "p", "", "Please enter a password")
	flag.StringVar(&code, "c", "", "Please enter an authorization code")
	flag.StringVar(&hostname, "h", "", "Please enter the hostname")
	flag.StringVar(&clientid, "ci", "", "Please enter a ClientID")
	flag.StringVar(&clientsecret, "cs", "", "Please enter a Client Secret")
	flag.StringVar(&search, "s", "", "Please enter a search string")
	flag.StringVar(&limit, "l", "100", "Please enter a search string")
	flag.Parse()

	// obtain the authorization code in browser: https://server/SASLogon/oauth/authorize?client_id=client&response_type=code
	ai := core.AuthInfo{
		// Username:     username,
		// Password:     password,
		Code:         code,
		GrantType:    "authorization_code",
		ClientID:     clientid,
		ClientSecret: clientsecret}

	baseURL := hostname
	token := ai.GetToken(baseURL)
	queryfl := url.Values{}
	queryfl.Add("limit", limit)
	if search != "" {
		queryfl.Add("filter", "contains(name, "+search+")")
	}
	querymem := url.Values{}
	querymem.Add("filter", "contains(contentType, file)")

	fl := sasobjs.GetFolders(baseURL, token, queryfl)
	for _, folder := range fl.Items {
		fmt.Printf("Id: %v Name: %v Members: %v\n", folder.ID, folder.Name, folder.MemberCount)
		mem := sasobjs.GetMembers(baseURL, folder.ID, token, querymem)
		for _, member := range mem.Items {
			fmt.Printf("Member Name: %s Member URI: %s Member ID: %s\n", member.Name, member.URI, member.ID)
			fmt.Println(">>>")
			fmt.Println(string(sasobjs.GetFileContent(baseURL, member.URI, token, nil)))
			fmt.Println(">>>")
		}
	}
}
