package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v38/github"
	"golang.org/x/oauth2"
)

func main() {
	// Attempt to read the token from the keystore
	token, err := readGitHubToken()
	if err != nil {
		fmt.Println(err.Error())
		// If the secret doesn't exist or there's an error reading it, prompt the user
		token, err = getGitHubTokenFromUser()
		if err != nil {
			fmt.Printf("Error reading user input: %v\n", err)
			os.Exit(1)
		}

		// Save the token to the keystore
		err = writeGitHubToken(token)
		if err != nil {
			fmt.Printf("Error writing token to keystore: %v\n", err)
			os.Exit(2)
		}
	}

	// Initialize a GitHub api client with the access token
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)
	api := &GitHubClient{client}

	orgRepos, err := api.GetOrgRepositories()
	if err != nil {
		fmt.Printf("Error fetching org repositories: %v\n", err)
		os.Exit(1)
	}

	starredRepos, err := api.GetStarredRepositories()
	if err != nil {
		fmt.Printf("Error fetching starred repositories: %v\n", err)
		os.Exit(1)
	}

	createBookmarksHTML(append(orgRepos, starredRepos...))
}
