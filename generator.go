package main

import (
        "context"
        "fmt"
        "os"
        "strings"
        "time"

        "github.com/google/go-github/v63/github"
        "github.com/joho/godotenv"
        "github.com/urfave/cli/v2"
)

func loadEnv() error {
        err := godotenv.Load()
        if err != nil {
                return fmt.Errorf("failed to load environment variables: %w", err)
        }
        return nil
}

func parseArgs() (string, string, string, string, time.Duration, int, error) {
        app := cli.NewApp()
        app.Name = "codeowners-generator"
        app.Usage = "Generates CODEOWNERS file based on top contributors"

        app.Flags = []cli.Flag{
                cli.StringFlag{
                        Name:    "github-server-url",
                        EnvVar:  "GITHUB_SERVER_URL",
                        Usage:   "GitHub server URL",
                        Value:   "https://api.github.com",
                },
                cli.StringFlag{
                        Name:    "organization-name",
                        EnvVar:  "ORGANIZATION_NAME",
                        Usage:   "GitHub organization name",
                },
                cli.StringFlag{
                        Name:    "repository-name",
                        EnvVar:  "REPOSITORY_NAME",
                        Usage:   "GitHub repository name",
                },
                cli.StringFlag{
                        Name:    "github-token",
                        EnvVar:  "GITHUB_TOKEN",
                        Usage:   "GitHub personal access token",
                },
                cli.DurationFlag{
                        Name:    "duration",
                        EnvVar:  "DURATION",
                        Usage:   "Time period for contributor analysis (e.g., 30d)",
                        Value:   time.Hour * 24 * 30, // Default to 30 days
                },
                cli.IntFlag{
                        Name:    "code-reviewers-count",
                        EnvVar:  "CODE_REVIEWERS_COUNT",
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

                return generateCodeowners(githubServerURL, organizationName, repositoryName, githubToken, duration, codeReviewersCount)
        }

        err := app.Run(os.Args)
        if err != nil {
                return "", "", "", "", time.Duration(0), 0, fmt.Errorf("failed to parse command line arguments: %w", err)
        }

        return githubServerURL, organizationName, repositoryName, githubToken, duration, codeReviewersCount, nil
}

func getGithubClient(githubServerURL, githubToken string) (*github.Client, error) {
        ts := github.TokenSource{Token: githubToken}
        tc := github.NewClient(github.ClientOptions{
                HTTPClient: &http.Client{Transport: http.DefaultTransport},
                BaseURL:    githubServerURL,
                UserAgent:  "codeowners-generator",
        })
        tc.Transport = http.DefaultTransport
        tc.RateLimit = nil
        tc.UserAgent = "codeowners-generator"
        return tc, nil
}

func getTopContributors(ctx context.Context, client *github.Client, owner, repo string, duration time.Duration) ([]github.User, error) {
        // Fetch commit activity for the specified duration
        since := time.Now().Add(-duration)
        commits, _, err := client.Activity.ListCommits(ctx, owner, repo, nil, &github.CommitsListOptions{
                Since: since.Format("2006-01-02T15:04:05Z"),
        })
        if err != nil {
                return nil, fmt.Errorf("failed to fetch commits: %w", err)
        }

        // Group contributors by login and count their commits
        contributorCounts := make(map[string]int)
        for _, commit := range commits {
                if commit.Author != nil {
                        contributorCounts[commit.Author.Login]++
                }
        }

        // Sort contributors by commit count in descending order
        contributors := make([]github.User, 0, len(contributorCounts))
        for login, count := range contributorCounts {
                contributors = append(contributors, github.User{Login: &login, Contributions: &count})
        }
        sort.Slice(contributors, func(i, j int) bool {
                return contributors[i].Contributions > contributors[j].Contributions
        })

        return contributors[:codeReviewersCount], nil
}

func generateCodeowners(githubServerURL, organizationName, repositoryName, githubToken string, duration time.Duration, codeReviewersCount int) error {
        ctx := context.Background()

        // Load environment variables
        if err := loadEnv(); err != nil {
                return fmt.Errorf("failed to load environment variables: %w", err)
        }

        // Create GitHub client
        client, err := getGithubClient(githubServerURL, githubToken)
        if err != nil {
                return fmt.Errorf("failed to create GitHub client: %w", err)
        }

        // Fetch top contributors
        topContributors, err := getTopContributors(ctx, client, organizationName, repositoryName, duration)
        if err != nil {
                return fmt.Errorf("failed to fetch top contributors: %w", err)
        }

        // Generate CODEOWNERS content
        var codeownersContent strings.Builder
        codeownersContent.WriteString("# This file is automatically generated by codeowners-generator\n")
        for _, contributor := range topContributors {
                codeownersContent.WriteString(fmt.Sprintf("%s/*\n\t@%s\n", contributor.Login, contributor.Login))
        }

        // Fetch existing CODEOWNERS
        existingCodeowners, _, err := client.Repositories.DownloadContents(ctx, organizationName, repositoryName, "CODEOWNERS", nil)
        if err != nil && !strings.Contains(err.Error(), "not found") {
                return fmt.Errorf("failed to fetch existing CODEOWNERS: %w", err)
        }

        // Compare generated CODEOWNERS with existing
        isDifferent := true
        if existingCodeowners != nil {
                existingContent := string(existingCodeowners)
                isDifferent = !strings.EqualFold(existingContent, codeownersContent.String())
        }

        // If different, create a new branch and submit a pull request
        if isDifferent {
                branchName := fmt.Sprintf("update-codeowners-%d", time.Now().Unix())
                ref := github.Reference{Ref: &branchName}
                _, _, err := client.Git.CreateRef(ctx, organizationName, repositoryName, "heads/"+branchName, ref)
                if err != nil {
                        return fmt.Errorf("failed to create new branch: %w", err)
                }

                // Create a new commit
                commitMessage := "Update CODEOWNERS based on top contributors"
                tree, err := client.Git.CreateTree(ctx, organizationName, repositoryName, "", &github.Tree{
                        Base:  nil,
                        Mode:  "100644",
                        Type:  "blob",
                        Path:  "CODEOWNERS",
                        Content: &codeownersContent.String(),
                })
                if err != nil {
                        return fmt.Errorf("failed to create new tree: %w", err)
                }

                parentCommit, _, err := client.Repositories.GetCommit(ctx, organizationName, repositoryName, "master")
                if err != nil {
                        return fmt.Errorf("failed to get parent commit: %w", err)
                }

                commit := &github.Commit{
                        Message: &commitMessage,
                        Tree:    tree,
                        Parent:  []string{parentCommit.SHA},
                }

                _, _, err = client.Git.CreateCommit