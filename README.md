## Codeowners Generator

**codeowners-generator** is a Go program that automatically generates a `CODEOWNERS` file for a GitHub repository based on the top contributors for a specified time period.

### Features

* Generates `CODEOWNERS` based on most active contributors within a time window.
* Supports command-line arguments and environment variables for configuration.
* Fetches existing `CODEOWNERS` and compares it with the generated one.
* Optionally creates a new branch and submits a pull request with the updated `CODEOWNERS` file (if differences exist).

### Installation

1. **Prerequisites:** Ensure you have Go installed ([https://golang.org/](https://golang.org/)).
2. **Clone the repository:**

```bash
git clone https://github.com/your-username/codeowners-generator.git
```

3. **Install dependencies:**

```bash
cd codeowners-generator
go mod download
```

### Usage

**1. Command-Line Arguments:**

```
codeowners-generator \
  -github-server-url=https://api.github.com \
  -organization-name=your-organization \
  -repository-name=your-repo \
  -github-token=YOUR_GITHUB_TOKEN \
  -duration=30d \
  -code-reviewers-count=3
```

* `-github-server-url`: Optional URL for the GitHub server (defaults to "[https://api.github.com](https://api.github.com)").
* `-organization-name`: Your GitHub organization name.
* `-repository-name`: The name of the repository.
* `-github-token`: Your personal access token with `repo` permission. (**Replace with your actual token, not shown here**)
* `-duration`: Time period for contributor analysis (e.g., `30d` for 30 days). Defaults to 30 days.
* `-code-reviewers-count`: Number of top contributors to include as code owners. Defaults to 3.

**2. Environment Variables:**

| Environment Variable | Description |
|---|---|
| `GITHUB_SERVER_URL` | GitHub server URL (defaults to `https://api.github.com`) |
| `ORGANIZATION_NAME` | GitHub organization name |
| `REPOSITORY_NAME` | GitHub repository name |
| `GITHUB_TOKEN` | GitHub personal access token with `repo` permission |
| `DURATION` | Time period for contributor analysis (e.g., `30d` for 30 days) |
| `CODE_REVIEWERS_COUNT` | Number of top contributors to include as code owners |

**3. Behavior:**

The program performs the following actions:

* Connects to the specified GitHub repository using your token.
* Fetches the top contributors based on commit activity within the provided duration.
* Generates a `CODEOWNERS` file assigning ownership to the top contributor usernames.
* Fetches the existing `CODEOWNERS` from the repository (if any).
* Compares the generated and existing files.
* If they differ, creates a new branch, commits the updated `CODEOWNERS`, and submits a pull request for review.

### Security

* **Protect your GitHub token:** Do not share your personal access token. Consider using environment variables or a secure storage mechanism.
* **Limited permissions:** Ensure your token has only the necessary permissions (e.g., `repo`).

### Contributing

We welcome contributions! Feel free to submit pull requests with improvements or bug fixes.

### License

This program is licensed under the MIT License.
