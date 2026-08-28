package github

import (
	"context"
	"fmt"
	"strings"

	ghErrors "github.com/github/github-mcp-server/pkg/errors"
	"github.com/github/github-mcp-server/pkg/inventory"
	"github.com/github/github-mcp-server/pkg/scopes"
	"github.com/github/github-mcp-server/pkg/translations"
	"github.com/github/github-mcp-server/pkg/utils"
	"github.com/google/go-github/v89/github"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// UpdateRepository creates a tool to update repository metadata, including renaming a repository.
func UpdateRepository(t translations.TranslationHelperFunc) inventory.ServerTool {
	return NewTool(
		ToolsetMetadataRepos,
		mcp.Tool{
			Name:        "update_repository",
			Description: t("TOOL_UPDATE_REPOSITORY_DESCRIPTION", "Update a GitHub repository, including renaming it or changing common repository metadata"),
			Annotations: &mcp.ToolAnnotations{
				Title:        t("TOOL_UPDATE_REPOSITORY_USER_TITLE", "Update repository"),
				ReadOnlyHint: false,
			},
			InputSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"owner": {
						Type:        "string",
						Description: "Repository owner (username or organization)",
					},
					"repo": {
						Type:        "string",
						Description: "Current repository name",
					},
					"name": {
						Type:        "string",
						Description: "New repository name. Providing this renames the repository.",
					},
					"description": {
						Type:        "string",
						Description: "New repository description. Use an empty string to clear it.",
					},
					"homepage": {
						Type:        "string",
						Description: "New repository homepage URL. Use an empty string to clear it.",
					},
					"private": {
						Type:        "boolean",
						Description: "Whether the repository should be private.",
					},
					"default_branch": {
						Type:        "string",
						Description: "New default branch name. The branch must already exist.",
					},
					"archived": {
						Type:        "boolean",
						Description: "Whether the repository should be archived.",
					},
					"has_issues": {
						Type:        "boolean",
						Description: "Whether GitHub Issues are enabled.",
					},
					"has_projects": {
						Type:        "boolean",
						Description: "Whether GitHub Projects are enabled.",
					},
					"has_wiki": {
						Type:        "boolean",
						Description: "Whether the repository wiki is enabled.",
					},
					"is_template": {
						Type:        "boolean",
						Description: "Whether the repository is available as a template repository.",
					},
					"allow_squash_merge": {
						Type:        "boolean",
						Description: "Whether squash merging pull requests is allowed.",
					},
					"allow_merge_commit": {
						Type:        "boolean",
						Description: "Whether merge commits are allowed for pull requests.",
					},
					"allow_rebase_merge": {
						Type:        "boolean",
						Description: "Whether rebase merging pull requests is allowed.",
					},
					"allow_auto_merge": {
						Type:        "boolean",
						Description: "Whether auto-merge is allowed.",
					},
					"delete_branch_on_merge": {
						Type:        "boolean",
						Description: "Whether head branches are automatically deleted after pull requests merge.",
					},
					"allow_update_branch": {
						Type:        "boolean",
						Description: "Whether pull request head branches can be updated when behind the base branch.",
					},
				},
				Required: []string{"owner", "repo"},
			},
		},
		scopes.RequireAll(scopes.Repo),
		func(ctx context.Context, deps ToolDependencies, _ *mcp.CallToolRequest, args map[string]any) (*mcp.CallToolResult, any, error) {
			owner, err := RequiredParam[string](args, "owner")
			if err != nil {
				return utils.NewToolResultError(err.Error()), nil, nil
			}
			repo, err := RequiredParam[string](args, "repo")
			if err != nil {
				return utils.NewToolResultError(err.Error()), nil, nil
			}

			update := &github.Repository{}
			changed := false

			if value, ok, err := OptionalParamOK[string](args, "name"); err != nil {
				return utils.NewToolResultError(err.Error()), nil, nil
			} else if ok {
				if strings.TrimSpace(value) == "" {
					return utils.NewToolResultError("parameter name must not be empty"), nil, nil
				}
				update.Name = github.Ptr(value)
				changed = true
			}
			if value, ok, err := OptionalParamOK[string](args, "description"); err != nil {
				return utils.NewToolResultError(err.Error()), nil, nil
			} else if ok {
				update.Description = github.Ptr(value)
				changed = true
			}
			if value, ok, err := OptionalParamOK[string](args, "homepage"); err != nil {
				return utils.NewToolResultError(err.Error()), nil, nil
			} else if ok {
				update.Homepage = github.Ptr(value)
				changed = true
			}
			if value, ok, err := OptionalParamOK[bool](args, "private"); err != nil {
				return utils.NewToolResultError(err.Error()), nil, nil
			} else if ok {
				update.Private = github.Ptr(value)
				changed = true
			}
			if value, ok, err := OptionalParamOK[string](args, "default_branch"); err != nil {
				return utils.NewToolResultError(err.Error()), nil, nil
			} else if ok {
				if strings.TrimSpace(value) == "" {
					return utils.NewToolResultError("parameter default_branch must not be empty"), nil, nil
				}
				update.DefaultBranch = github.Ptr(value)
				changed = true
			}
			if value, ok, err := OptionalParamOK[bool](args, "archived"); err != nil {
				return utils.NewToolResultError(err.Error()), nil, nil
			} else if ok {
				update.Archived = github.Ptr(value)
				changed = true
			}
			if value, ok, err := OptionalParamOK[bool](args, "has_issues"); err != nil {
				return utils.NewToolResultError(err.Error()), nil, nil
			} else if ok {
				update.HasIssues = github.Ptr(value)
				changed = true
			}
			if value, ok, err := OptionalParamOK[bool](args, "has_projects"); err != nil {
				return utils.NewToolResultError(err.Error()), nil, nil
			} else if ok {
				update.HasProjects = github.Ptr(value)
				changed = true
			}
			if value, ok, err := OptionalParamOK[bool](args, "has_wiki"); err != nil {
				return utils.NewToolResultError(err.Error()), nil, nil
			} else if ok {
				update.HasWiki = github.Ptr(value)
				changed = true
			}
			if value, ok, err := OptionalParamOK[bool](args, "is_template"); err != nil {
				return utils.NewToolResultError(err.Error()), nil, nil
			} else if ok {
				update.IsTemplate = github.Ptr(value)
				changed = true
			}
			if value, ok, err := OptionalParamOK[bool](args, "allow_squash_merge"); err != nil {
				return utils.NewToolResultError(err.Error()), nil, nil
			} else if ok {
				update.AllowSquashMerge = github.Ptr(value)
				changed = true
			}
			if value, ok, err := OptionalParamOK[bool](args, "allow_merge_commit"); err != nil {
				return utils.NewToolResultError(err.Error()), nil, nil
			} else if ok {
				update.AllowMergeCommit = github.Ptr(value)
				changed = true
			}
			if value, ok, err := OptionalParamOK[bool](args, "allow_rebase_merge"); err != nil {
				return utils.NewToolResultError(err.Error()), nil, nil
			} else if ok {
				update.AllowRebaseMerge = github.Ptr(value)
				changed = true
			}
			if value, ok, err := OptionalParamOK[bool](args, "allow_auto_merge"); err != nil {
				return utils.NewToolResultError(err.Error()), nil, nil
			} else if ok {
				update.AllowAutoMerge = github.Ptr(value)
				changed = true
			}
			if value, ok, err := OptionalParamOK[bool](args, "delete_branch_on_merge"); err != nil {
				return utils.NewToolResultError(err.Error()), nil, nil
			} else if ok {
				update.DeleteBranchOnMerge = github.Ptr(value)
				changed = true
			}
			if value, ok, err := OptionalParamOK[bool](args, "allow_update_branch"); err != nil {
				return utils.NewToolResultError(err.Error()), nil, nil
			} else if ok {
				update.AllowUpdateBranch = github.Ptr(value)
				changed = true
			}

			if !changed {
				return utils.NewToolResultError("no repository updates were provided"), nil, nil
			}

			client, err := deps.GetClient(ctx)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to get GitHub client: %w", err)
			}

			updated, resp, err := client.Repositories.Edit(ctx, owner, repo, update)
			if err != nil {
				return ghErrors.NewGitHubAPIErrorResponse(ctx, "failed to update repository", resp, err), nil, nil
			}
			if resp != nil && resp.Body != nil {
				defer func() { _ = resp.Body.Close() }()
			}

			result := struct {
				ID            int64  `json:"id"`
				Name          string `json:"name"`
				FullName      string `json:"full_name"`
				URL           string `json:"html_url"`
				Description   string `json:"description,omitempty"`
				Private       bool   `json:"private"`
				DefaultBranch string `json:"default_branch,omitempty"`
				Archived      bool   `json:"archived"`
			}{
				ID:            updated.GetID(),
				Name:          updated.GetName(),
				FullName:      updated.GetFullName(),
				URL:           updated.GetHTMLURL(),
				Description:   updated.GetDescription(),
				Private:       updated.GetPrivate(),
				DefaultBranch: updated.GetDefaultBranch(),
				Archived:      updated.GetArchived(),
			}

			return MarshalledTextResult(result), nil, nil
		},
	)
}
