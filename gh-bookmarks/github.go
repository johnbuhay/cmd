package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v38/github"
	"golang.org/x/crypto/ssh/terminal"
)

// GitHubClient embeds the github.Client so I can create this "class"
type GitHubClient struct {
	*github.Client
}

// GetGitHubOrg checks args for input and defaults to the username associated with the GitHub token if none are provided
func (gh *GitHubClient) GetGitHubOrg() string {
	org := ""
	if len(os.Args) > 1 {
		org = os.Args[1]
	}

	if len(os.Args) > 2 {
		fmt.Println("warn: extra arguments found, only using the first")
	}

	if org == "" {
		// If no input is provided, default to the username associated with the GitHub token
		username, _ := gh.GetGitHubUsername()
		org = username
	}

	return org
}

func readGitHubToken() (string, error) {
	secret, err := getToken()
	if err != nil {
		return string(secret), err
	}
	return strings.TrimSpace(string(secret)), nil
}

func writeGitHubToken(input string) error {
	return storeToken(input)
}

// GetGitHubUsername returns the username associated with the provided token, otherwise returns with error
func (gh *GitHubClient) GetGitHubUsername() (string, error) {
	user, _, err := gh.Client.Users.Get(context.Background(), "")
	if err != nil {
		return "", err
	}
	return *user.Login, nil
}

// getGitHubTokenFromUser reads input from the user without echoing it to the terminal
func getGitHubTokenFromUser() (string, error) {
	fmt.Print("GitHub access token not found. Please enter your GitHub access token: ")
	bytePassword, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()

	if err != nil {
		return "", err
	}
	p := string(bytePassword)
	return p, nil
}

func (gh *GitHubClient) GetOrgRepositories() ([]*github.Repository, error) {
	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	// Set the org value with a function that prompts the user
	// default retrieves the GitHub username using the token
	org := gh.GetGitHubOrg()
	var repos []*github.Repository
	for {
		orgRepos, resp, err := gh.Client.Repositories.List(context.Background(), org, opt)
		if err != nil {
			fmt.Printf("Error fetching repositories: %v\n", err)
			os.Exit(1)
		}
		repos = append(repos, orgRepos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return repos, nil
}

// GetStarredRepositories returns a list of repositories that are starred by the user
// but not of the StarredRepository type
func (gh *GitHubClient) GetStarredRepositories() ([]*github.Repository, error) {
	opt := &github.ActivityListStarredOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var repos []*github.Repository
	for {
		starredRepos, resp, err := gh.Client.Activity.ListStarred(context.Background(), "", opt)
		if err != nil {
			fmt.Printf("Error fetching repositories: %v\n", err)
			os.Exit(1)
		}
		// https://pkg.go.dev/github.com/google/go-github/v38/github#StarredRepository
		for _, repo := range starredRepos {
			repos = append(repos, repo.Repository)
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return repos, nil
}
