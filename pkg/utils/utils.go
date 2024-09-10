package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

func ParseArgs() (string, string, string, string, time.Duration, int, error) {
	app := cli.NewApp()
	app.Name = "codeowners-generator"
	app.Usage = "Generates CODEOWNERS file based on top contributors"

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "github-server-url",
			EnvVars: []string{"GITHUB_SERVER_URL"},
			Usage:   "GitHub server URL",
			Value:   "https://api.github.com",
		},
		&cli.StringFlag{
			Name:    "organization-name",
			EnvVars: []string{"ORGANIZATION_NAME"},
			Usage:   "GitHub organization name",
		},
		&cli.StringFlag{
			Name:    "repository-name",
			EnvVars: []string{"REPOSITORY_NAME"},
			Usage:   "GitHub repository name",
		},
		&cli.StringFlag{
			Name:    "github-token",
			EnvVars: []string{"GITHUB_TOKEN"},
			Usage:   "GitHub personal access token",
		},
		&cli.DurationFlag{
			Name:    "duration",
			EnvVars: []string{"DURATION"},
			Usage:   "Time period for contributor analysis (e.g., 30d)",
			Value:   time.Hour * 24 * 30, // Default to 30 days
		},
		&cli.IntFlag{
			Name:    "code-reviewers-count",
			EnvVars: []string{"CODE_REVIEWERS_COUNT"},
			Usage:   "Number of top contributors to include as code owners",
			Value:   3, // Default to 3
		},
	}

	app.Action = func(c *cli.Context) error {
		githubServerURL := c.String("github-server-url")
		organizationName := c.String("organization-name")
		repositoryName := c.String("repository-name")
		githubToken := c.String("github-token")
		duration := c.Duration("duration")
		codeReviewersCount := c.Int("code-reviewers-count")

		fmt.Println()

		return nil // Handle error from main function
	}

	err := app.Run(os.Args)
	if err != nil {
		return "", "", "", "", time.Duration(0), 0, fmt.Errorf("failed to parse command line arguments: %w", err)
	}

	return githubServerURL, organizationName, repositoryName, githubToken, duration, codeReviewersCount, nil
}
