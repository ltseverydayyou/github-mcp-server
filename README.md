# GitHub MCP Server

[![GitHub contributors](https://img.shields.io/github/contributors/github/github-mcp-server)](https://github.com/github/github-mcp-server/graphs/contributors)
[![GitHub stars](https://img.shields.io/github/stars/github/github-mcp-server)](https://github.com/github/github-mcp-server/stargazers)
[![GitHub license](https://img.shields.io/github/license/github/github-mcp-server)](https://github.com/github/github-mcp-server/blob/main/LICENSE)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/github/github-mcp-server)

The GitHub MCP Server connects AI tools directly to GitHub's platform. It enables AI tools, agents and assistants to read repositories and code, manage issues and pull requests, analyze code, monitor activity, and automate workflows—all through natural language interactions.

### Use Cases

- **Repository Management**: Browse and query code, search files, analyze commits, and understand project structure across any repository you have access to.
- **Issue & Pull Request Automation**: Create, update, and manage issues and pull requests. Let your AI assist with triaging bugs, reviewing code, and maintaining project workflows.
- **CI/CD & Workflow Intelligence**: Monitor GitHub Actions workflow runs, investigate failures, manage releases, and gain insights into your development pipeline.
- **Code Analysis**: Inspect code, examine security findings, analyze Dependabot alerts, understand code patterns, and get detailed insights into your codebase.
- **Team Collaboration**: Access discussions, manage notifications, analyze team activity, streamline processes, and stay connected with your development team.

Designed for developers who want to connect their AI tools to GitHub context and capabilities, from simple natural language queries to complex multi-step workflows.

---

## Remote GitHub MCP Server

[![Install in VS Code](https://img.shields.io/badge/VS_Code-Install_Server-0098FF?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=github&config=%7B%22type%22%3A%22http%22%2C%22url%22%3A%22https%3A%2F%2Fapi.githubcopilot.com%2Fmcp%2F%22%7D) [![Install in VS Code Insiders](https://img.shields.io/badge/VS_Code_Insiders-Install_Server-24bfa5?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=github&config=%7B%22type%22%3A%22http%22%2C%22url%22%3A%2F%2Fapi.githubcopilot.com%2Fmcp%2F%22%7D&quality=insiders)

The remote GitHub MCP server is hosted by GitHub and is the easiest way to get up and running. If your MCP host does not support remote MCP servers, you can use the [local GitHub MCP server](#local-github-mcp-server) instead.

### Prerequisites

1. A compatible MCP host with remote server support (VS Code 1.101+, Copilot in other IDEs, Claude Desktop, Claude Code, Cursor, Windsurf, etc.)
2. Any relevant [policies enabled](https://docs.github.com/en/copilot/how-tos/administer-copilot/manage-for-organization/manage-policies#configuring-mcp-server-access) 

### Configure in VS Code 

For more information about configuring MCP servers in VS Code, see the [official VS Code documentation](https://code.visualstudio.com/docs/copilot/chat/mcp-servers).

#### Remote Server with OAuth

[![Install in VS Code](https://img.shields.io/badge/VS_Code-Install_Server-0098FF?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=github&config=%7B%22type%22%3A%22http%22%2C%22url%22%3A%22https%3A%2F%2Fapi.githubcopilot.com%2Fmcp%2F%22%7D) [![Install in VS Code Insiders](https://img.shields.io/badge/VS_Code_Insiders-Install_Server-24bfa5?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=github&config=%7B%22type%22%3A%22http%22%2C%22url%22%3A%22https%3A%2F%2Fapi.githubcopilot.com%2Fmcp%2F%22%7D&quality=insiders)

> **Note:** The remote GitHub MCP server supports OAuth authorization headers by default. MCP clients can request the server's protected resource metadata, discover the GitHub authorization server, dynamically register as OAuth clients, and launch a browser-based authorization flow.

For clients that support MCP Apps, the server can also return interactive UI resources when the client advertises the `io.modelcontextprotocol/ui` extension during initialization. Use the `github-mcp-server` tool `generate-app-csp` to print the required CSP entries for supported UI bundles.

### GitHub Enterprise Cloud with data residency (ghe.com)

Remote hosting currently supports `https://api.githubcopilot.com` for github.com. GitHub Enterprise Cloud with data residency (`ghe.com`) requires the local MCP server. Configure a local server with `GITHUB_HOST=https://<subdomain>.ghe.com` as described below.

### Remote Server without OAuth

To configure the remote server without OAuth, you can use a GitHub Personal Access Token (PAT). We recommend using a [fine-grained PAT](https://github.com/settings/personal-access-tokens/new) where you can scope your token to specific repositories and permissions. For detailed guidance on creating and securing your PAT, see [Creating a personal access token](https://github.com/github/github-mcp-server/blob/main/docs/installation-guides/install-copilot-cli.md#creating-a-pat) in the GitHub Copilot CLI installation guide.

#### Using headers

```json
{
  "servers": {
    "github": {
      "type": "http",
      "url": "https://api.githubcopilot.com/mcp/",
      "headers": {
        "Authorization": "Bearer ${input:github_mcp_pat}"
      }
    }
  },
  "inputs": [
    {
      "type": "promptString",
      "id": "github_mcp_pat",
      "description": "GitHub Personal Access Token",
      "password": true
    }
  ]
}
```

#### Using `mcp-remote`

If your MCP host does not natively support remote MCP servers or OAuth, you can use [mcp-remote](https://www.npmjs.com/package/mcp-remote) to connect to the remote GitHub MCP server using a PAT.

```json
{
  "servers": {
    "github": {
      "command": "npx",
      "args": [
        "-y",
        "mcp-remote",
        "https://api.githubcopilot.com/mcp/",
        "--header",
        "Authorization: Bearer ${input:github_mcp_pat}"
      ]
    }
  },
  "inputs": [
    {
      "type": "promptString",
      "id": "github_mcp_pat",
      "description": "GitHub Personal Access Token",
      "password": true
    }
  ]
}
```

### Configure in other MCP hosts

- **[Claude Code](docs/installation-guides/install-claude.md)** - Setup guide for Claude Code (CLI and VS Code extension)
- **[Claude Desktop](docs/installation-guides/install-claude-desktop.md)** - Setup guide for Claude Desktop
- **[Codex](docs/installation-guides/install-codex.md)** - Setup guide for OpenAI Codex CLI and IDE extension
- **[Cursor](docs/installation-guides/install-cursor.md)** - Setup guide for Cursor IDE
- **[Gemini CLI](docs/installation-guides/install-gemini-cli.md)** - Setup guide for Gemini CLI
- **[GitHub Copilot CLI](docs/installation-guides/install-copilot-cli.md)** - Setup guide for GitHub Copilot CLI
- **[GitHub Copilot in VS Code](docs/installation-guides/install-vscode.md)** - Setup guide for GitHub Copilot in VS Code
- **[Visual Studio](docs/installation-guides/install-visual-studio.md)** - Setup guide for Visual Studio
- **[Windsurf](docs/installation-guides/install-windsurf.md)** - Setup guide for Windsurf

### Remote Server Configuration

For information about configuring the remote server via headers, toolsets, and individual tools, see [Remote Server Configuration](docs/remote-server.md). Remote-only toolsets and tools are also documented there.

---

## Local GitHub MCP Server

[![Install with Docker in VS Code](https://img.shields.io/badge/VS_Code-Install_Server-0098FF?style=flat-square&logo=docker&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=github&inputs=%5B%7B%22id%22%3A%22github_token%22%2C%22type%22%3A%22promptString%22%2C%22description%22%3A%22GitHub%20Personal%20Access%20Token%22%2C%22password%22%3Atrue%7D%5D&config=%7B%22command%22%3A%22docker%22%2C%22args%22%3A%5B%22run%22%2C%22-i%22%2C%22--rm%22%2C%22-e%22%2C%22GITHUB_PERSONAL_ACCESS_TOKEN%22%2C%22ghcr.io%2Fgithub%2Fgithub-mcp-server%22%5D%2C%22env%22%3A%7B%22GITHUB_PERSONAL_ACCESS_TOKEN%22%3A%22%24%7Binput%3Agithub_token%7D%22%7D%7D) [![Install with Docker in VS Code Insiders](https://insiders.vscode.dev/redirect/mcp/install?name=github&inputs=%5B%7B%22id%22%3A%22github_token%22%2C%22type%22%3A%22promptString%22%2C%22description%22%3A%22GitHub%20Personal%20Access%20Token%22%2C%22password%22%3Atrue%7D%5D&config=%7B%22command%22%3A%22docker%22%2C%22args%22%3A%5B%22run%22%2C%22-i%22%2C%22--rm%22%2C%22-e%22%2C%22GITHUB_PERSONAL_ACCESS_TOKEN%22%2C%22ghcr.io%2Fgithub%2Fgithub-mcp-server%22%5D%2C%22env%22%3A%7B%22GITHUB_PERSONAL_ACCESS_TOKEN%22%3A%22%24%7Binput%3Agithub_token%7D%22%7D%7D&quality=insiders)

You can use the GitHub MCP server locally via Docker or build from source.

### Prerequisites

1. To use the server in a container, install [Docker](https://www.docker.com/).
2. After Docker is installed, ensure Docker is running. If pulling the Docker image fails, you may need to sign out of `ghcr.io` with `docker logout ghcr.io` before retrying.
3. Create a [GitHub Personal Access Token](https://github.com/settings/personal-access-tokens/new). The token should have the permissions you need for the repositories you want to access.

### Docker usage

For remote transports with the Docker image, run the server with port `8080` exposed:

```bash
docker run --rm -i -p 8080:8080 \
  -e GITHUB_PERSONAL_ACCESS_TOKEN \
  ghcr.io/github/github-mcp-server http
```

For stdio mode:

```bash
docker run --rm -i \
  -e GITHUB_PERSONAL_ACCESS_TOKEN \
  ghcr.io/github/github-mcp-server stdio
```

### Build from source

```bash
go build ./cmd/github-mcp-server
```

### Configuration

The local server supports the same toolsets and tools documented below. Use `--toolsets` to select toolsets and `--tools` to select individual tools.

### Toolsets

The following sets of tools are available:

<details>
<summary><strong>Toolsets</strong></summary>

<!-- GENERATED: TOOLSETS -->

- **actions**
  - `actions_get`
  - `actions_list`
  - `actions_run_trigger`
  - `get_job_logs`
- **code_quality**
  - `get_code_quality_finding`
- **code_security**
  - `get_code_scanning_alert`
  - `list_code_scanning_alerts`
- **context**
  - `get_me`
  - `get_team_members`
  - `get_teams`
- **copilot**
  - `assign_copilot_to_issue`
  - `request_copilot_review`
- **copilot_issue_intents**
  - `assign_copilot_to_issue_with_intent`
- **dependabot**
  - `get_dependabot_alert`
  - `list_dependabot_alerts`
- **discussions**
  - `discussion_comment_write`
  - `get_discussion`
  - `get_discussion_comments`
  - `list_discussion_categories`
  - `list_discussions`
- **git**
  - `get_repository_tree`
- **gists**
  - `create_gist`
  - `get_gist`
  - `list_gists`
  - `update_gist`
- **issues**
  - `add_issue_comment`
  - `find_duplicate`
  - `issue_dependency_read`
  - `issue_dependency_write`
  - `issue_read`
  - `issue_write`
  - `list_issue_fields`
  - `list_issue_types`
  - `list_issues`
  - `search_issues`
- **labels**
  - `get_label`
  - `label_write`
  - `list_label`
- **notifications**
  - `dismiss_notification`
  - `get_notification_details`
  - `list_notifications`
  - `manage_notification_subscription`
  - `manage_repository_notification_subscription`
  - `mark_all_notifications_read`
- **orgs**
  - `search_orgs`
- **projects**
  - `projects_get`
  - `projects_list`
  - `projects_write`
- **pull_requests**
  - `add_comment_to_pending_review`
  - `add_reply_to_pull_request_comment`
  - `create_pull_request`
  - `list_pull_requests`
  - `merge_pull_request`
  - `pull_request_read`
  - `pull_request_review_write`
  - `request_copilot_review`
  - `search_pull_requests`
  - `update_pull_request`
  - `update_pull_request_branch`
- **repos**
  - `create_branch`
  - `create_or_update_file`
  - `create_repository`
  - `delete_file`
  - `delete_repository`
  - `fork_repository`
  - `get_commit`
  - `get_file_blame`
  - `get_file_contents`
  - `get_latest_release`
  - `get_release_by_tag`
  - `get_tag`
  - `github_api`
  - `github_file_download`
  - `github_file_upload`
  - `github_graphql`
  - `github_repository_file_upload`
  - `list_branches`
  - `list_commits`
  - `list_releases`
  - `list_repository_collaborators`
  - `list_starred_repositories`
  - `list_tags`
  - `push_files`
  - `release_write`
  - `search_code`
  - `search_commits`
  - `search_repositories`
  - `star_repository`
  - `unstar_repository`
  - `update_repository`
- **secret_protection**
  - `get_secret_scanning_alert`
  - `list_secret_scanning_alerts`
- **security_advisories**
  - `get_global_security_advisory`
  - `list_global_security_advisories`
  - `list_org_repository_security_advisories`
  - `list_repository_security_advisories`
- **stargazers**
  - `list_starred_repositories`
  - `star_repository`
  - `unstar_repository`
- **users**
  - `search_users`

<!-- END GENERATED: TOOLSETS -->

</details>

### Tools

<details>
<summary><strong>Tools</strong></summary>

<!-- GENERATED: TOOLS -->

- **actions_get** - Get Actions resource
  - `method`: The method to execute (string, required)
  - `owner`: Repository owner (string, required)
  - `repo`: Repository name (string, required)
  - `resource_id`: The unique identifier of the resource. This will vary based on the "method" provided, so ensure you provide the correct ID:
- Provide a workflow ID or workflow file name (e.g. ci.yaml) for 'get_workflow' method.
- Provide a workflow run ID for 'get_workflow_run', 'get_workflow_run_usage', and 'get_workflow_run_logs_url' methods.
- Provide an artifact ID for 'download_workflow_run_artifact' method.
- Provide a job ID for 'get_workflow_job' method. (string, required)

- **actions_list** - List Actions resources
  - `method`: The action to perform (string, required)
  - `owner`: Repository owner (string, required)
  - `page`: Page number for pagination (default: 1) (number, optional)
  - `per_page`: Results per page for pagination (default: 30, max: 100) (number, optional)
  - `repo`: Repository name (string, required)
  - `resource_id`: The unique identifier of the resource. This will vary based on the "method" provided, so ensure you provide the correct ID:
- Do not provide any resource ID for 'list_workflows' method.
- Provide a workflow ID or workflow file name (e.g. ci.yaml) for 'list_workflow_runs' method, or omit to list all workflow runs.
- Provide a workflow run ID for 'list_workflow_jobs' and 'list_workflow_run_artifacts' methods. (string, optional)
  - `workflow_jobs_filter`: Filters for workflow jobs. **ONLY** used when method is 'list_workflow_jobs' (object, optional)
  - `workflow_runs_filter`: Filters for workflow runs. **ONLY** used when method is 'list_workflow_runs' (object, optional)

- **actions_run_trigger** - Trigger Actions workflow operation
  - `inputs`: Inputs the workflow accepts. Only used for 'run_workflow' method. (object, optional)
  - `method`: The method to execute (string, required)
  - `owner`: Repository owner (string, required)
  - `ref`: The git reference for the workflow. The reference can be a branch or tag name. Required for 'run_workflow' method. (string, optional)
  - `repo`: Repository name (string, required)
  - `run_id`: The ID of the workflow run. Required for all methods except 'run_workflow'. (number, optional)
  - `workflow_id`: The workflow ID (numeric) or workflow file name (e.g. main.yml, ci.yaml). Required for 'run_workflow' method. (string, optional)

- **add_comment_to_pending_review** - Add comment to pending review
  - **OAuth Challenge Scopes**: `repo`
  - `body`: The text of the review comment (string, required)
  - `line`: The line of the blob in the pull request diff that the comment applies to. For multi-line comments, the last line of the range (number, optional)
  - `owner`: Repository owner (string, required)
  - `path`: The relative path to the file that necessitates a comment (string, required)
  - `pullNumber`: Pull request number (number, required)
  - `repo`: Repository name (string, required)
  - `side`: The side of the diff to comment on. LEFT indicates the previous state, RIGHT indicates the new state (string, optional)
  - `startLine`: For multi-line comments, the first line of the range that the comment applies to (number, optional)
  - `startSide`: For multi-line comments, the starting side of the diff that the comment applies to. LEFT indicates the previous state, RIGHT indicates the new state (string, optional)
  - `subjectType`: The level at which the comment is targeted (string, required)

- **add_issue_comment** - Add issue comment or reaction
  - **OAuth Challenge Scopes**: `repo`
  - `body`: Comment content. Required unless reaction is provided. (string, optional)
  - `comment_id`: The numeric ID of the issue or pull request comment to react to. Use this for reactions to comments; omit it to react to the issue or pull request itself. Cannot be combined with body. (number, optional)
  - `issue_number`: Issue or pull request number to comment on or react to. (number, required)
  - `owner`: Repository owner (string, required)
  - `reaction`: Emoji reaction to add. Required unless body is provided. (string, optional)
  - `repo`: Repository name (string, required)

- **add_reply_to_pull_request_comment** - Reply or react to pull request comment
  - **OAuth Challenge Scopes**: `repo`
  - `body`: The text of the reply. Required unless reaction is provided. (string, optional)
  - `commentId`: The numeric ID of the pull request review comment to reply or react to. Use the number from a #discussion_r... anchor, not the GraphQL thread node ID (PRRT_...). (number, required)
  - `owner`: Repository owner (string, required)
  - `pullNumber`: Pull request number. Required when body is provided. (number, optional)
  - `reaction`: Emoji reaction to add. Required unless body is provided. (string, optional)
  - `repo`: Repository name (string, required)

- **assign_copilot_to_issue** - Assign Copilot to an issue
  - **OAuth Challenge Scopes**: `repo`
  - `base_ref`: Git reference (e.g., branch) that the agent will start its work from. If not specified, defaults to the repository's default branch (string, optional)
  - `custom_instructions`: Optional custom instructions to guide the agent beyond the issue body. Use this to provide additional context, constraints, or guidance that is not captured in the issue description (string, optional)
  - `issue_number`: Issue number (number, required)
  - `owner`: Repository owner (string, required)
  - `repo`: Repository name (string, required)

- **assign_copilot_to_issue_with_intent** - Assign Copilot to issue with intent metadata
  - **OAuth Challenge Scopes**: `repo`
  - `base_ref`: Git reference (e.g., branch) that the agent will start its work from. If not specified, defaults to the repository's default branch. Ignored when is_suggestion is true (string, optional)
  - `confidence`: How confident you are in this choice. 'HIGH' for clear signal or explicit user request, 'MEDIUM' for reasonable inference with some ambiguity, 'LOW' for best guess with limited signal. (string, required)
  - `custom_instructions`: Optional custom instructions to guide the agent beyond the issue body. Ignored when is_suggestion is true (string, optional)
  - `is_suggestion`: If true, records a pending Copilot assignment intent rather than launching the agent. Approval later supplies the launch context; base_ref and custom_instructions are ignored in this case. (boolean, required)
  - `issue_number`: Issue number (number, required)
  - `owner`: Repository owner (string, required)
  - `rationale`: One concise sentence explaining what specifically about the issue led to choosing Copilot. State the concrete signal (e.g. 'Well-scoped task with clear acceptance criteria'). (string, required)
  - `repo`: Repository name (string, required)

- **create_branch** - Create branch
  - **OAuth Challenge Scopes**: `repo`
  - `branch`: Name for new branch (string, required)
  - `from_branch`: Source branch (defaults to repo default) (string, optional)
  - `owner`: Repository owner (string, required)
  - `repo`: Repository name (string, required)

- **create_gist** - Create gist
  - `content`: Content for simple single-file gist creation (string, required)
  - `description`: Description of the gist (string, optional)
  - `filename`: Filename for simple single-file gist creation (string, required)
  - `public`: Whether the gist is public (boolean, optional)

- **create_or_update_file** - Create or update file
  - **OAuth Challenge Scopes**: `repo`
  - `allow_symlink_write`: Set true to update a symbolic link itself; content must be its new target path. (boolean, optional)
  - `branch`: Branch to create/update the file in (string, required)
  - `content`: Content of the file, exactly as it should appear once written. Do not base64-encode it; this server does that before calling the REST API. (string, required)
  - `message`: Commit message (string, required)
  - `owner`: Repository owner (username or organization) (string, required)
  - `path`: Path where to create/update the file (string, required)
  - `repo`: Repository name (string, required)
  - `sha`: The blob SHA of the file being replaced. Required if the file already exists. (string, optional)

- **create_pull_request** - Create pull request
  - **OAuth Challenge Scopes**: `repo`
  - `base`: Branch to merge into (string, required)
  - `body`: PR description (string, optional)
  - `draft`: Create as draft PR (boolean, optional)
  - `head`: Branch containing changes (string, required)
  - `maintainer_can_modify`: Allow maintainer edits (boolean, optional)
  - `owner`: Repository owner (string, required)
  - `repo`: Repository name (string, required)
  - `reviewers`: GitHub usernames or ORG/team-slug team reviewers to request reviews from (string[], optional)
  - `title`: PR title (string, required)

- **create_repository** - Create repository
  - `autoInit`: Initialize with README (boolean, optional)
  - `description`: Repository description (string, optional)
  - `name`: Repository name (string, required)
  - `organization`: Organization to create the repository in (omit to create in your personal account) (string, optional)
  - `private`: Whether the repository should be private. Defaults to true (private) when omitted. (boolean, optional)

- **delete_file** - Delete file
  - **OAuth Challenge Scopes**: `repo`
  - `branch`: Branch to delete the file from (string, required)
  - `message`: Commit message (string, required)
  - `owner`: Repository owner (username or organization) (string, required)
  - `path`: Path to the file to delete (string, required)
  - `repo`: Repository name (string, required)

- **delete_repository** - Delete repository
  - **OAuth Challenge Scopes**: `delete_repo`
  - `owner`: Repository owner (string, required)
  - `repo`: Repository name (string, required)

- **discussion_comment_write** - Write discussion comment
  - **OAuth Challenge Scopes**: `repo`
  - `body`: Comment content (required for 'add', 'reply', and 'update' methods) (string, optional)
  - `commentNodeID`: The Node ID of the discussion comment (required for 'reply', 'update', 'delete', 'mark_answer', and 'unmark_answer' methods). For 'reply', this is the top-level comment to reply to; GitHub Discussions only support one level of nesting. (string, optional)
  - `discussionNumber`: Discussion number (required for 'add' and 'reply' methods) (number, optional)
  - `method`: Write operation to perform on a discussion comment.
Options are:
- 'add' - adds a new top-level comment to a discussion.
- 'reply' - replies to a top-level discussion comment (GitHub Discussions only support one level of nesting).
- 'update' - updates an existing discussion comment.
- 'delete' - deletes a discussion comment.
- 'mark_answer' - marks a discussion comment as the answer (Q&A only).
- 'unmark_answer' - unmarks a discussion comment as the answer (Q&A only). (string, required)
  - `owner`: Repository owner (required for 'add' and 'reply' methods) (string, optional)
  - `repo`: Repository name (required for 'add' and 'reply' methods) (string, optional)

- **dismiss_notification** - Dismiss notification
  - `state`: The new state of the notification (read/done) (string, required)
  - `threadID`: The ID of the notification thread (string, required)

- **find_duplicate** - Find duplicate issue
  - `issue_number`: The issue number to find duplicates for (number, required)
  - `owner`: Repository owner (string, required)
  - `repo`: Repository name (string, required)

- **fork_repository** - Fork repository
  - **OAuth Challenge Scopes**: `repo`
  - `organization`: Organization to fork to (string, optional)
  - `owner`: Repository owner (string, required)
  - `repo`: Repository name (string, required)

- **get_code_quality_finding** - Get code quality finding
  - **OAuth Challenge Scopes**: `repo`
  - `findingNumber`: The number of the finding. (number, required)
  - `owner`: The owner of the repository. (string, required)
  - `repo`: The name of the repository. (string, required)

- **get_code_scanning_alert** - Get code scanning alert
  - **OAuth Challenge Scopes**: `security_events`
  - `alertNumber`: The number of the alert. (number, required)
  - `owner`: The owner of the repository. (string, required)
  - `repo`: The name of the repository. (string, required)

- **get_commit** - Get commit details
  - **OAuth Challenge Scopes**: `repo`
  - `detail`: Level of detail to include for changed files. "none" omits stats and files entirely. "stats" (default) includes per-file metadata: filename, status, and lines-of-code counts (additions, deletions, changes), with no patch content. "full_patch" additionally includes the unified diff content for each file and can be very large. (string, optional)
  - `owner`: Repository owner (string, required)
  - `page`: Page number for pagination (min 1) (number, optional)
  - `perPage`: Results per page for pagination (min 1, max 100) (number, optional)
  - `repo`: Repository name (string, required)
  - `sha`: Commit SHA, branch name, or tag name (string, required)

- **get_dependabot_alert** - Get Dependabot alert
  - **OAuth Challenge Scopes**: `repo`
  - `alertNumber`: The number of the alert. (number, required)
  - `owner`: The owner of the repository. (string, required)
  - `repo`: The name of the repository. (string, required)

- **get_discussion** - Get discussion
  - **OAuth Challenge Scopes**: `repo`
  - `discussionNumber`: Discussion Number (number, required)
  - `owner`: Repository owner (string, required)
  - `repo`: Repository name (string, required)

- **get_discussion_comments** - Get discussion comments
  - **OAuth Challenge Scopes**: `repo`
  - `after`: Cursor for pagination. Use the cursor from the previous response. (string, optional)
  - `discussionNumber`: Discussion Number (number, required)
  - `includeReplies`: When true, each top-level comment will include its replies nested within it (up to 100 replies per comment, which is the GitHub API maximum). Defaults to false. (boolean, optional)
  - `owner`: Repository owner (string, required)
  - `perPage`: Results per page for pagination (min 1, max 100) (number, optional)
  - `repo`: Repository name (string, required)

- **get_file_blame** - Get file blame information
  - **OAuth Challenge Scopes**: `repo`
  - `owner`: Repository owner (string, required)
  - `path`: Path to file (string, required)
  - `ref`: The git ref to fetch blame for (branch, tag, SHA). If not provided, the repository's default branch is used. (string, optional)
  - `repo`: Repository name (string, required)
  - `sha`: Optional commit SHA. If provided, it takes precedence over `ref`. (string, optional)

- **get_file_contents** - Get file or directory contents
  - **OAuth Challenge Scopes**: `repo`
  - `fields`: Subset of fields to return for each entry when the path is a directory. If omitted, all fields are returned. Ignored when the path is a single file. Use this to reduce response size when listing directories and you only need specific fields, e.g. just 'name' and 'type'. (string[], optional)
  - `owner`: Repository owner (username or organization) (string, required)
  - `path`: Path to file/directory (string, optional)
  - `ref`: Accepts optional git refs such as `refs/tags/{tag}`, `refs/heads/{branch}` or `refs/pull/{pr_number}/head` (string, optional)
  - `repo`: Repository name (string, required)
  - `sha`: Accepts optional commit SHA. If specified, it will be used instead of ref (string, optional)

- **get_global_security_advisory** - Get a global security advisory
  - `ghsaId`: GitHub Security Advisory ID (format: GHSA-xxxx-xxxx-xxxx). (string, required)

- **get_job_logs** - Get job logs
  - **OAuth Challenge Scopes**: `repo`
  - `failed_only`: When true, gets logs for all failed jobs in the workflow run specified by run_id. Requires run_id to be provided. (boolean, optional)
  - `job_id`: The unique identifier of the workflow job. Required when getting logs for a single job. (number, optional)
  - `owner`: Repository owner (string, required)
  - `repo`: Repository name (string, required)
  - `return_content`: Returns actual log content instead of URLs (boolean, optional)
  - `run_id`: The unique identifier of the workflow run. Required when failed_only is true to get logs for all failed jobs in a workflow run. (number, optional)
  - `tail_lines`: Number of lines to return from the end of the log (number, optional)

- **get_label** - Get label
  - **OAuth Challenge Scopes**: `repo`
  - `name`: Label name. (string, required)
  - `owner`: Repository owner (username or organization name) (string, required)
  - `repo`: Repository name (string, required)

- **get_latest_release** - Get latest release
  - **OAuth Challenge Scopes**: `repo`
  - `owner`: Repository owner (string, required)
  - `repo`: Repository name (string, required)

- **get_me** - Get current user
  - `userOnly`: If true, returns only the current user data without GitHub context. (boolean, optional)

- **get_notification_details** - Get notification details
  - **OAuth Challenge Scopes**: `notifications`
  - `notificationID`: The ID of the notification (string, required)

- **get_release_by_tag** - Get a release by tag name
  - **OAuth Challenge Scopes**: `repo`
  - `owner`: Repository owner (string, required)
  - `repo`: Repository name (string, required)
  - `tag`: Tag name (e.g., 'v1.0.0') (string, required)

- **get_tag** - Get tag details
  - **OAuth Challenge Scopes**: `repo`
  - `owner`: Repository owner (string, required)
  - `repo`: Repository name (string, required)
  - `tag`: Tag name (string, required)

- **github_api** - Call GitHub REST API
  - **OAuth Challenge Scopes**: `repo`
  - `accept`: Optional Accept header for GitHub API previews or alternate representations. (string, optional)
  - `body`: Optional JSON request body. (object, optional)
  - `endpoint`: Relative GitHub REST API endpoint without a scheme or host. (string, required)
  - `method`: HTTP method. (string, required)

- **github_file_download** - Download complete repository file
  - **OAuth Challenge Scopes**: `repo`
  - `owner`: Repository owner. (string, required)
  - `path`: Repository file path to download. (string, required)
  - `ref`: Optional git ref such as a branch, tag, or refs/pull/<number>/head. (string, optional)
  - `repo`: Repository name. (string, required)
  - `sha`: Optional commit SHA. Takes precedence over ref. (string, optional)

- **github_file_upload** - Upload local file to GitHub
  - **OAuth Challenge Scopes**: `repo`
  - `content_type`: Alias for media_type; takes precedence when both are supplied. (string, optional)
  - `endpoint`: Relative GitHub upload endpoint without a scheme or host. Existing query parameters are preserved. (string, required)
  - `file_path`: Local filesystem path on the MCP server host. Remote URLs are rejected. (string, required)
  - `filename`: Alias for name. Used when name is omitted. (string, optional)
  - `label`: Optional label query parameter. (string, optional)
  - `media_type`: Optional MIME/media type. Defaults to application/octet-stream. (string, optional)
  - `name`: Optional filename/name query parameter. Defaults to the file basename. (string, optional)
  - `query`: Optional additional query parameters appended to the endpoint. (object, optional)

- **github_graphql** - Call GitHub GraphQL API
  - **OAuth Challenge Scopes**: `repo`
  - `query`: GraphQL query or mutation document. (string, required)
  - `variables`: Optional GraphQL variables object. (object, optional)

- **github_repository_file_upload** - Commit ChatGPT file to repository
  - **OAuth Challenge Scopes**: `repo`
  - `branch`: Branch to update. Omit to use the repository default branch. (string, optional)
  - `file`:  (object, required)
  - `message`: Commit message. (string, required)
  - `owner`: Repository owner. (string, required)
  - `path`: Repository path to create or replace. (string, required)
  - `repo`: Repository name. (string, required)
  - `sha`: Optional expected blob SHA for an existing target. When omitted, the current SHA is detected automatically. (string, optional)

- **list_branches** - List branches
  - **OAuth Challenge Scopes**: `repo`
  - `owner`: Repository owner (string, required)
  - `page`: Page number for pagination (min 1) (number, optional)
  - `perPage`: Results per page for pagination (min 1, max 100) (number, optional)
  - `repo`: Repository name (string, required)

- **list_commits** - List commits
  - **OAuth Challenge Scopes**: `repo`
  - `author`: Author username or email address to filter commits by (string, optional)
  - `fields`: Subset of fields to return for each commit. If omitted, all fields are returned. Use this to reduce response size when you only need specific fields, e.g. just 'sha' and 'html_url'. (string[], optional)
  - `owner`: Repository owner (string, required)
  - `page`: Page number for pagination (min 1) (number, optional)
  - `path`: Only commits containing this file path will be returned (string, optional)
  - `perPage`: Results per page for pagination (min 1, max 100) (number, optional)
  - `repo`: Repository name (string, required)
  - `sha`: Commit SHA, branch or tag name to list commits of. If not provided, uses the default branch of the repository. If a commit SHA is provided, will list commits up to that SHA. (string, optional)
  - `since`: Only commits after this date will be returned (ISO 8601 format: YYYY-MM-DDTHH:MM:SSZ or YYYY-MM-DD) (string, optional)
  - `until`: Only commits before this date (ISO 8601 timestamp) (string, optional)

- **list_releases** - List releases
  - **OAuth Challenge Scopes**: `repo`
  - `fields`: Subset of fields to return for each release. If omitted, all fields are returned. Use this to reduce response size when you only need specific fields; omitting 'body' in particular drops the largest per-release data. (string[], optional)
  - `owner`: Repository owner (string, required)
  - `page`: Page number for pagination (min 1) (number, optional)
  - `perPage`: Results per page for pagination (min 1, max 100) (number, optional)
  - `repo`: Repository name (string, required)

- **list_repository_collaborators** - List repository collaborators
  - **OAuth Challenge Scopes**: `repo`
  - `affiliation`: Filter by affiliation. Can be one of: 'outside' (outside collaborators), 'direct' (all with permissions regardless of org membership), 'all' (all collaborators). Default: 'all' (string, optional)
  - `owner`: Repository owner (string, required)
  - `page`: Page number for pagination (default 1, min 1) (number, optional)
  - `perPage`: Results per page for pagination (default 30, min 1, max 100) (number, optional)
  - `repo`: Repository name (string, required)

- **list_starred_repositories** - List starred repositories
  - `direction`: The direction to sort the results by. (string, optional)
  - `page`: Page number for pagination (min 1) (number, optional)
  - `perPage`: Results per page for pagination (min 1, max 100) (number, optional)
  - `sort`: How to sort the results. Can be either 'created' (when the repository was starred) or 'updated' (when the repository was last pushed to). (string, optional)
  - `username`: Username to list starred repositories for. Defaults to the authenticated user. (string, optional)

- **list_tags** - List tags
  - **OAuth Challenge Scopes**: `repo`
  - `owner`: Repository owner (string, required)
  - `page`: Page number for pagination (min 1) (number, optional)
  - `perPage`: Results per page for pagination (min 1, max 100) (number, optional)
  - `repo`: Repository name (string, required)

- **manage_notification_subscription** - Manage notification subscription
  - **OAuth Challenge Scopes**: `notifications`
  - `action`: Action to perform: ignore, watch, or delete the notification subscription. (string, required)
  - `notificationID`: The ID of the notification thread. (string, required)

- **manage_repository_notification_subscription** - Manage repository notification subscription
  - **OAuth Challenge Scopes**: `notifications`
  - `action`: Action to perform: ignore, watch, or delete repository notifications subscription for the provided repository. (string, required)
  - `owner`: The account owner of the repository. (string, required)
  - `repo`: The name of the repository. (string, required)

- **mark_all_notifications_read** - Mark all notifications read
  - **OAuth Challenge Scopes**: `notifications`
  - `lastReadAt`: Describes the last point that notifications were checked (optional). Default: Now (string, optional)
  - `owner`: Optional repository owner. If provided with repo, only notifications for this repository are marked as read. (string, optional)
  - `repo`: Optional repository name. If provided with owner, only notifications for this repository are marked as read. (string, optional)

- **merge_pull_request** - Merge pull request
  - **OAuth Challenge Scopes**: `repo`
  - `commit_message`: Extra detail for merge commit (string, optional)
  - `commit_title`: Title for merge commit (string, optional)
  - `merge_method`: Merge method (string, optional)
  - `owner`: Repository owner (string, required)
  - `pullNumber`: Pull request number (number, required)
  - `repo`: Repository name (string, required)

- **projects_get** - Get Projects resource
  - **OAuth Challenge Scopes**: `repo`, `project`, `read:project`
  - `field_id`: The field's ID. Required for 'get_project_field' method. (number, optional)
  - `field_names`: Specific list of field names to include in the response when getting a project item (e.g. ["Status", "Priority"]). Resolved server-side to field IDs — pass this instead of 'fields' when you only know the human-readable names. Mutually exclusive with 'fields' — provide one, not both. Only used for 'get_project_item' method. (string[], optional)
  - `fields`: Specific list of field IDs to include in the response when getting a project item (e.g. ["102589", "985201", "169875"]). If neither 'fields' nor 'field_names' is provided, only the title field is included. Mutually exclusive with 'field_names' — provide one, not both. Only used for 'get_project_item' method. (string[], optional)
  - `item_id`: The item's ID. Required for 'get_project_item' method. (number, optional)
  - `method`: The method to execute (string, required)
  - `owner`: The owner (user or organization login). The name is not case sensitive. (string, optional)
  - `owner_type`: Owner type (user or org). If not provided, will be automatically detected. (string, optional)
  - `project_number`: The project's number. (number, optional)
  - `status_update_id`: The node ID of the project status update. Required for 'get_project_status_update' method. (string, optional)
  - `view_id`: The node ID of the project view. Required for 'get_project_view' method. (string, optional)

- **projects_list** - List Projects resources
  - **OAuth Challenge Scopes**: `repo`, `project`, `read:project`
  - `after`: Forward pagination cursor from previous pageInfo.nextCursor. (string, optional)
  - `before`: Backward pagination cursor from previous pageInfo.prevCursor (rare). (string, optional)
  - `field_names`: Field names to include when listing project items (e.g. ["Status", "Priority"]). Resolved server-side to field IDs — pass this instead of 'fields' when you only know the human-readable names. Names that fail to resolve return a structured error. Mutually exclusive with 'fields' — provide one, not both. Only used for 'list_project_items' method. (string[], optional)
  - `fields`: Field IDs to include when listing project items (e.g. ["102589", "985201"]). CRITICAL: Always provide to get field values. Without this (and without 'field_names'), only titles returned. Mutually exclusive with 'field_names' — provide one, not both. Only used for 'list_project_items' method. (string[], optional)
  - `method`: The action to perform (string, required)
  - `owner`: The owner (user or organization login). The name is not case sensitive. (string, required)
  - `owner_type`: Owner type (user or org). If not provided, will automatically try both. (string, optional)
  - `per_page`: Results per page (max 50) (number, optional)
  - `project_number`: The project's number. Required for 'list_project_fields', 'list_project_items', 'list_project_views', and 'list_project_status_updates' methods. (number, optional)
  - `query`: Filter/query string. For list_projects: filter by title text and state (e.g. "roadmap is:open"). For list_project_items: advanced filtering using GitHub's project filtering syntax. (string, optional)

- **projects_write** - Write Projects resource
  - **OAuth Challenge Scopes**: `repo`, `project`
  - `body`: The body of the status update (markdown). Used for 'create_project_status_update' method. (string, optional)
  - `field_name`: The name of the iteration field (e.g. 'Sprint'). Required for 'create_iteration_field' method. (string, optional)
  - `filter`: Saved view filter; omit on update to preserve it, or pass null to clear it. (string or null, optional)
  - `issue_number`: The issue number. Required for 'add_project_item' when item_type is 'issue'. Also accepted by 'update_project_item' to resolve the item by issue number (combine with item_owner and item_repo). (number, optional)
  - `item_id`: The project item ID. Required for 'delete_project_item'. For 'update_project_item', provide either item_id, or (item_owner + item_repo + issue_number) to resolve the item by issue. (number, optional)
  - `item_owner`: The owner (user or organization) of the repository containing the issue or pull request. Required for 'add_project_item' method. Also accepted by 'update_project_item' when resolving the item by issue number. (string, optional)
  - `item_repo`: The name of the repository containing the issue or pull request. Required for 'add_project_item' method. Also accepted by 'update_project_item' when resolving the item by issue number. (string, optional)
  - `item_type`: The item's type, either issue or pull_request. Required for 'add_project_item' method. (string, optional)
  - `items`: The items to update with the top-level 'updated_field'. Required for 'update_project_items'; prefer it over calling 'update_project_item' in a loop. Each entry must match exactly one reference variant: 'node_id', numeric 'item_id', or 'item_owner' + 'item_repo' + 'issue_number'. Limit: 50 items per call. (object[], optional)
  - `iteration_duration`: Duration in days for iterations of the field (e.g. 7 for weekly, 14 for bi-weekly). Required for 'create_iteration_field' method. (number, optional)
  - `iterations`: Custom iterations for 'create_iteration_field' method. Only set this when you need iterations with varying durations, breaks between them, or specific titles. Otherwise omit it: GitHub auto-creates three iterations of 'iteration_duration' days starting on 'start_date', which is the right choice for most cases. (object[], optional)
  - `layout`: View layout; required when creating a view. (string, optional)
  - `method`: The method to execute (string, required)
  - `name`: View name; required when creating a view. (string, optional)
  - `owner`: The project owner (user or organization login). The name is not case sensitive. (string, required)
  - `owner_type`: Owner type (user or org). Required for 'create_project' method. If not provided for other methods, will be automatically detected. (string, optional)
  - `project_number`: The project's number. Required for all methods except 'create_project'. (number, optional)
  - `pull_request_number`: The pull request number (use when item_type is 'pull_request' for 'add_project_item' method). Provide either issue_number or pull_request_number. (number, optional)
  - `start_date`: Start date in YYYY-MM-DD format. Used for 'create_project_status_update' and 'create_iteration_field' methods. (string, optional)
  - `status`: The status of the project. Used for 'create_project_status_update' method. (string, optional)
  - `target_date`: The target date of the status update in YYYY-MM-DD format. Used for 'create_project_status_update' method. (string, optional)
  - `title`: The project title. Required for 'create_project' method. (string, optional)
  - `updated_field`: The field/value to apply, using {"id": 123, "value": ...} or {"name": "Status", "value": ...}; null clears the field. Required for 'update_project_item' and 'update_project_items', where one top-level field/value applies to every item in a batch. For 'update_project_item' SINGLE_SELECT fields, the name form accepts option names; the ID form expects an option ID. (object or null, optional)
  - `view_id`: Project view node ID for update or delete; must belong to owner/project_number. (string, optional)
  - `visible_field_names`: Ordered project field names to show on create or replace on update; omit on update to preserve, or pass [] to reset. Mutually exclusive with visible_fields. Roadmap accepts only []. (string[], optional)
  - `visible_fields`: Ordered project field database IDs to show on create or replace on update; omit on update to preserve, or pass [] to reset. Mutually exclusive with visible_field_names. Roadmap accepts only []. (string[], optional)

- **pull_request_read** - Read pull request
  - **OAuth Challenge Scopes**: `repo`
  - `after`: Cursor for pagination, used only by the get_review_comments method. Pass the endCursor from the previous page's PageInfo to fetch the next page. (string, optional)
  - `method`: Action to specify what pull request data needs to be retrieved from GitHub.
Possible options:
1. get - Get details of a specific pull request.
2. get_diff - Get the diff of a pull request.
3. get_status - Get combined commit status of a head commit in a pull request.
4. get_files - Get the list of files changed in a pull request. Use with pagination parameters to control the number of results returned.
5. get_commits - Get the list of commits on a pull request. Use with pagination parameters to control the number of results returned.
6. get_review_comments - Get review threads on a pull request. Each thread contains logically grouped review comments made on the same code location during pull request reviews. Returns threads with metadata (isResolved, isOutdated, isCollapsed) and their associated comments. Use cursor-based pagination (perPage, after) to control results.
7. get_reviews - Get the reviews on a pull request. When asked for review comments, use get_review_comments method. Use with pagination parameters to control results.
8. get_comments - Get comments on a pull request. Use this if user doesn't specifically want review comments. Use with pagination parameters to control results.
9. get_check_runs - Get check runs for the head commit of a pull request. Check runs are the individual CI/CD jobs and checks that run on the PR. (string, required)
  - `owner`: Repository owner (string, required)
  - `page`: Page number for pagination (min 1) (number, optional)
  - `perPage`: Results per page for pagination (min 1, max 100) (number, optional)
  - `pullNumber`: Pull request number (number, required)
  - `repo`: Repository name (string, required)

- **pull_request_review_write** - Write pull request review
  - **OAuth Challenge Scopes**: `repo`
  - `body`: Review comment text (string, optional)
  - `commitID`: SHA of commit to review (string, optional)
  - `event`: Review action to perform. (string, optional)
  - `method`: The write operation to perform on pull request review. (string, required)
  - `owner`: Repository owner (string, required)
  - `pullNumber`: Pull request number (number, required)
  - `repo`: Repository name (string, required)
  - `threadId`: The node ID of the review thread (e.g. PRRT_kwDOxxx). Required for resolve_thread and unresolve_thread methods. Get thread IDs from pull_request_read with method get_review_comments. (string, optional)

- **push_files** - Push multiple files
  - **OAuth Challenge Scopes**: `repo`
  - `branch`: Branch to push to (string, required)
  - `files`: Array of file objects to push, each object with path (string) and content (string) (object[], required)
  - `message`: Commit message (string, required)
  - `owner`: Repository owner (string, required)
  - `repo`: Repository name (string, required)

- **release_write** - Write release
  - **OAuth Challenge Scopes**: `repo`
  - `asset_id`: Release asset ID for update_asset or delete_asset. (number, optional)
  - `body`: Release notes/body. (string, optional)
  - `content_base64`: Base64-encoded bytes for upload_asset. Mutually exclusive with file_path. (string, optional)
  - `discussion_category_name`: Optional discussion category to create for the release. (string, optional)
  - `draft`: Whether the release is a draft. (boolean, optional)
  - `file_path`: Local filesystem path on the MCP server host for upload_asset. Mutually exclusive with content_base64; streamed directly when provided. (string, optional)
  - `generate_release_notes`: Ask GitHub to generate release notes when creating the release. (boolean, optional)
  - `label`: Optional release asset label. (string, optional)
  - `make_latest`: Latest-release behavior. (string, optional)
  - `media_type`: Optional MIME type for upload_asset. (string, optional)
  - `method`: Release operation to perform. (string, required)
  - `name`: Release name; for asset operations, the asset filename/name. (string, optional)
  - `owner`: Repository owner. (string, required)
  - `prerelease`: Whether the release is a prerelease. (boolean, optional)
  - `release_id`: Release ID for update, delete, or upload_asset. (number, optional)
  - `repo`: Repository name. (string, required)
  - `state`: Optional asset state for update_asset. (string, optional)
  - `tag_name`: Tag name. Required for create; optional for update. (string, optional)
  - `target_commitish`: Target branch or commitish. (string, optional)

- **request_copilot_review** - Request Copilot review
  - **OAuth Challenge Scopes**: `repo`
  - `owner`: Repository owner (string, required)
  - `pullNumber`: Pull request number (number, required)
  - `repo`: Repository name (string, required)

- **search_code** - Search code
  - **OAuth Challenge Scopes**: `repo`
  - `fields`: Subset of fields to return for each code search result. If omitted, all fields are returned. Use this to reduce response size when you only need specific fields; omitting 'repository' and 'text_matches' in particular drops the largest per-result data. (string[], optional)
  - `order`: Sort order for results (string, optional)
  - `page`: Page number for pagination (min 1) (number, optional)
  - `perPage`: Results per page for pagination (min 1, max 100) (number, optional)
  - `query`: Search query (GitHub code search REST). Implicit AND between terms; supports `OR`, `NOT`, and `"quoted phrase"` for exact match. Qualifiers: `repo:owner/repo`, `org:`, `user:`, `language:`, `path:dir` (prefix match), `filename:exact.ext`, `extension:`, `in:file`, `in:path`, `size:`, `is:archived`, `is:fork`. Max 256 chars. Examples: `WithContext language:go org:github`; `"package main" repo:o/r`; `func extension:go path:cmd repo:o/r`; `NOT TODO language:go repo:o/r`. (string, required)
  - `sort`: Sort field ('indexed' only) (string, optional)

- **search_commits** - Search commits
  - **OAuth Challenge Scopes**: `repo`
  - `order`: Sort order (string, optional)
  - `page`: Page number for pagination (min 1) (number, optional)
  - `perPage`: Results per page for pagination (min 1, max 100) (number, optional)
  - `query`: Commit search query (GitHub commit search REST). Searches commit messages on the default branch only. Scope the search with `repo:owner/repo`, `org:`, or `user:` (queries without a scope qualifier match across all of GitHub and are usually not what you want). Other qualifiers: `author:`, `committer:`, `author-name:`, `committer-name:`, `author-email:`, `committer-email:`, `author-date:`, `committer-date:` (supports `>`, `<`, `>=`, `<=`, and `YYYY-MM-DD..YYYY-MM-DD` ranges), `merge:true|false`, `hash:`, `tree:`, `parent:`, `is:public`. Examples: `repo:owner/repo fix panic`; `org:github author:defunkt committer-date:>=2024-01-01`; `"refactor cache" repo:o/r`; `hash:abc1234 repo:o/r`. (string, required)
  - `sort`: Sort by author or committer date (defaults to best match) (string, optional)

- **search_issues** - Search issues
  - **OAuth Challenge Scopes**: `repo`
  - `fields`: Subset of fields to return for each issue result. If omitted, all fields are returned. Use this to reduce response size when you only need specific fields; omitting 'body', 'reactions', and 'labels' in particular drops the largest per-result data. (string[], optional)
  - `order`: Sort order (string, optional)
  - `owner`: Optional repository owner. If provided with repo, only issues for this repository are listed. (string, optional)
  - `page`: Page number for pagination (min 1) (number, optional)
  - `perPage`: Results per page for pagination (min 1, max 100) (number, optional)
  - `query`: The search query, as natural language. When the user gives alternative wordings, include them as plain words rather than joining them with OR. (string, required)
  - `repo`: Optional repository name. If provided with owner, only issues for this repository are listed. (string, optional)
  - `sort`: Sort field by number of matches of categories, defaults to best match (string, optional)

- **search_orgs** - Search organizations
  - `order`: Sort order (string, optional)
  - `page`: Page number for pagination (min 1) (number, optional)
  - `perPage`: Results per page for pagination (min 1, max 100) (number, optional)
  - `query`: Organization search query. Examples: 'microsoft', 'location:california', 'created:>=2025-01-01'. Search is automatically scoped to type:org. (string, required)
  - `sort`: Sort field by category (string, optional)

- **search_pull_requests** - Search pull requests
  - **OAuth Challenge Scopes**: `repo`
  - `fields`: Subset of fields to return for each pull request result. If omitted, all fields are returned. Use this to reduce response size when you only need specific fields; omitting 'body', 'reactions', and 'labels' in particular drops the largest per-result data. (string[], optional)
  - `order`: Sort order (string, optional)
  - `owner`: Optional repository owner. If provided with repo, only pull requests for this repository are listed. (string, optional)
  - `page`: Page number for pagination (min 1) (number, optional)
  - `perPage`: Results per page for pagination (min 1, max 100) (number, optional)
  - `query`: Search query using GitHub pull request search syntax (string, required)
  - `repo`: Optional repository name. If provided with owner, only pull requests for this repository are listed. (string, optional)
  - `sort`: Sort field by number of matches of categories, defaults to best match (string, optional)

- **search_repositories** - Search repositories
  - `minimal_output`: Return minimal repository information (default: true). When false, returns full GitHub API repository objects. (boolean, optional)
  - `order`: Sort order (string, optional)
  - `page`: Page number for pagination (min 1) (number, optional)
  - `perPage`: Results per page for pagination (min 1, max 100) (number, optional)
  - `query`: Repository search query. Examples: 'machine learning in:name stars:>1000 language:python', 'topic:react', 'user:facebook'. Supports advanced search syntax for precise filtering. (string, required)
  - `sort`: Sort repositories by field, defaults to best match (string, optional)

- **search_users** - Search users
  - `order`: Sort order (string, optional)
  - `page`: Page number for pagination (min 1) (number, optional)
  - `perPage`: Results per page for pagination (min 1, max 100) (number, optional)
  - `query`: User search query. Examples: 'john smith', 'location:seattle', 'followers:>100'. Search is automatically scoped to type:user. (string, required)
  - `sort`: Sort users by number of followers or repositories, or when the person joined GitHub. (string, optional)

- **star_repository** - Star repository
  - `owner`: Repository owner (string, required)
  - `repo`: Repository name (string, required)

- **sub_issue_write** - Write sub-issue
  - **OAuth Challenge Scopes**: `repo`
  - `after_id`: The ID of the sub-issue to be prioritized after (either after_id OR before_id should be specified) (number, optional)
  - `before_id`: The ID of the sub-issue to be prioritized before (either after_id OR before_id should be specified) (number, optional)
  - `issue_number`: The number of the parent issue (number, required)
  - `method`: The action to perform on a single sub-issue
Options are:
- 'add' - add a sub-issue to a parent issue in a GitHub repository.
- 'remove' - remove a sub-issue from a parent issue in a GitHub repository.
- 'reprioritize' - change the order of sub-issues within a parent issue in a GitHub repository. Use either 'after_id' or 'before_id' to specify the new position.
Writes issue hierarchy. To move a sub-issue to a new parent, use `add` with `replace_parent=true`; there is no writable parent field. (string, required)
  - `owner`: Repository owner (string, required)
  - `replace_parent`: When true, replaces the sub-issue's current parent issue. Use with 'add' method only. (boolean, optional)
  - `repo`: Repository name (string, required)
  - `sub_issue_id`: The ID of the sub-issue to add. ID is not the same as issue number (number, required)

- **unstar_repository** - Unstar repository
  - `owner`: Repository owner (string, required)
  - `repo`: Repository name (string, required)

- **update_gist** - Update gist
  - `content`: Content for the file (string, required)
  - `description`: Updated description of the gist (string, optional)
  - `filename`: Filename to update or create (string, required)
  - `gist_id`: ID of the gist to update (string, required)

- **update_pull_request** - Update pull request
  - **OAuth Challenge Scopes**: `repo`
  - `base`: New base branch name (string, optional)
  - `body`: New description (string, optional)
  - `draft`: Mark pull request as draft (true) or ready for review (false) (boolean, optional)
  - `maintainer_can_modify`: Allow maintainer edits (boolean, optional)
  - `owner`: Repository owner (string, required)
  - `pullNumber`: Pull request number to update (number, required)
  - `repo`: Repository name (string, required)
  - `reviewers`: GitHub usernames or ORG/team-slug team reviewers to request reviews from (string[], optional)
  - `state`: New state (string, optional)
  - `title`: New title (string, optional)

- **update_pull_request_branch** - Update pull request branch
  - **OAuth Challenge Scopes**: `repo`
  - `expectedHeadSha`: The expected SHA of the pull request's HEAD ref (string, optional)
  - `owner`: Repository owner (string, required)
  - `pullNumber`: Pull request number (number, required)
  - `repo`: Repository name (string, required)

<!-- END GENERATED: TOOLS -->

</details>

## Resources

The server exposes repository content as MCP resources using URI templates. These resources let clients fetch file and directory contents without calling a tool directly.

## License

This project is licensed under the terms of the MIT open source license. Please refer to [LICENSE](LICENSE) for the full terms.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for details.
