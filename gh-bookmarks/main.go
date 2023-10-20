package main

import (
	"context"
	"fmt"
	"os"
	"os/user"

	"github.com/google/go-github/v38/github"
	"golang.org/x/oauth2"
)

func main() {
	// Get the user's home directory
	usr, err := user.Current()
	if err != nil {
		fmt.Printf("Error getting user's home directory: %v\n", err)
		os.Exit(1)
	}

	// Define the path to the token file
	tokenPath := usr.HomeDir + "/.local/.github-token"

	// Attempt to read the token from the file
	token, err := readGitHubToken(tokenPath)
	if err != nil {
		// If the file doesn't exist or there's an error reading it, prompt the user
		fmt.Println("GitHub access token not found. Please enter your GitHub access token:")
		_, err := fmt.Scan(&token)
		if err != nil {
			fmt.Printf("Error reading user input: %v\n", err)
			os.Exit(1)
		}

		// Save the token to the file
		err = writeGitHubToken(tokenPath, token)
		if err != nil {
			fmt.Printf("Error writing token to file: %v\n", err)
			os.Exit(1)
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
