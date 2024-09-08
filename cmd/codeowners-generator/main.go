package main

import (
	"context"
	"fmt"
	"os"

	"github.com/your-username/codeowners-generator/pkg/github"
	"github.com/your-username/codeowners-generator/pkg/utils"
)

func main() {
	// Parse command-line arguments
	githubServerURL, organizationName, repositoryName, githubToken, duration, codeReviewersCount, err := utils.ParseArgs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Create GitHub client
	client, err := github.NewClient(githubServerURL, githubToken)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Fetch top contributors
	ctx := context.Background()
	topContributors, err := github.GetTopContributors(ctx, client, organizationName, repositoryName, duration)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Generate CODEOWNERS content
	codeownersContent := github.GenerateCodeowners(topContributors)

	// Fetch existing CODEOWNERS
	existingCodeowners, err := github.FetchExistingCodeowners(ctx, client, organizationName, repositoryName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Compare generated CODEOWNERS with existing
	isDifferent := github.CompareCodeowners(codeownersContent, existingCodeowners)

	// If different, create a new branch and submit a pull request
	if isDifferent {
		err = github.CreateBranchAndSubmitPullRequest(ctx, client, organizationName, repositoryName, codeownersContent)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Pull request created successfully!")
	} else {
		fmt.Println("CODEOWNERS file is up-to-date.")
	}
}
