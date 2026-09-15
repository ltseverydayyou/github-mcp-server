package github

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
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

func generateReposToolsetInstructions(_ *inventory.Inventory) string {
	return `## Repository files

When a file already exists in the ChatGPT sandbox (normally /mnt/data), prefer github_repository_file_upload to create or replace it in a repository without copying the full contents into MCP arguments. This is especially important for large files. Do not create a GitHub Actions workflow merely to transfer or update a sandbox file when github_repository_file_upload can perform the update directly.

When complete repository file contents are needed locally, especially for files too large or awkward for get_file_contents, prefer github_file_download to place the full file in /mnt/data for inspection.`
}

// RepositoryFileUpload commits a ChatGPT sandbox file directly into a repository without sending its contents through MCP arguments.
func RepositoryFileUpload(t translations.TranslationHelperFunc) inventory.ServerTool {
	return NewTool(
		ToolsetMetadataRepos,
		mcp.Tool{
			Name:        "github_repository_file_upload",
			Description: t("TOOL_GITHUB_REPOSITORY_FILE_UPLOAD_DESCRIPTION", "Create or replace a repository file from a file in the ChatGPT sandbox. Reads the source directly from /mnt/data (or GITHUB_MCP_SANDBOX_ROOT), automatically resolves the existing blob SHA when needed, and avoids sending large file contents through MCP arguments."),
			Annotations: &mcp.ToolAnnotations{
				Title:           t("TOOL_GITHUB_REPOSITORY_FILE_UPLOAD_USER_TITLE", "Commit sandbox file to repository"),
				ReadOnlyHint:    false,
				DestructiveHint: jsonschema.Ptr(true),
			},
			InputSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"owner":     {Type: "string", Description: "Repository owner."},
					"repo":      {Type: "string", Description: "Repository name."},
					"path":      {Type: "string", Description: "Repository path to create or replace."},
					"file_path": {Type: "string", Description: "Source file in the ChatGPT sandbox, normally /mnt/data/<file>."},
					"message":   {Type: "string", Description: "Commit message."},
					"branch":    {Type: "string", Description: "Branch to update. Omit to use the repository default branch."},
					"sha":       {Type: "string", Description: "Optional expected blob SHA for an existing target. When omitted, the tool detects the current SHA automatically."},
				},
				Required: []string{"owner", "repo", "path", "file_path", "message"},
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
			repoPath, err := RequiredParam[string](args, "path")
			if err != nil {
				return utils.NewToolResultError(err.Error()), nil, nil
			}
			filePath, err := RequiredParam[string](args, "file_path")
			if err != nil {
				return utils.NewToolResultError(err.Error()), nil, nil
			}
			message, err := RequiredParam[string](args, "message")
			if err != nil {
				return utils.NewToolResultError(err.Error()), nil, nil
			}
			if strings.TrimSpace(repoPath) == "" || strings.HasSuffix(repoPath, "/") {
				return utils.NewToolResultError("path must point to a repository file"), nil, nil
			}
			if strings.TrimSpace(message) == "" {
				return utils.NewToolResultError("message must not be empty"), nil, nil
			}

			filePath, err = resolveSandboxPath(filePath, true)
			if err != nil {
				return utils.NewToolResultError(err.Error()), nil, nil
			}
			content, err := os.ReadFile(filePath)
			if err != nil {
				return utils.NewToolResultErrorFromErr("failed to read sandbox file", err), nil, nil
			}

			branch, _ := OptionalParam[string](args, "branch")
			targetSHA, _ := OptionalParam[string](args, "sha")
			client, err := deps.GetClient(ctx)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to get GitHub client: %w", err)
			}

			if targetSHA == "" {
				getOpts := &github.RepositoryContentGetOptions{Ref: branch}
				existing, directory, resp, getErr := client.Repositories.GetContents(ctx, owner, repo, repoPath, getOpts)
				if getErr == nil {
					defer closeGitHubResponse(resp)
					if existing == nil || directory != nil {
						return utils.NewToolResultError("path resolves to a directory, not a file"), nil, nil
					}
					targetSHA = existing.GetSHA()
				} else if resp == nil || resp.StatusCode != 404 {
					return ghErrors.NewGitHubAPIErrorResponse(ctx, "failed to inspect repository target", resp, getErr), nil, nil
				} else {
					closeGitHubResponse(resp)
				}
			}

			opts := &github.RepositoryContentFileOptions{
				Message: github.Ptr(message),
				Content: content,
			}
			if branch != "" {
				opts.Branch = github.Ptr(branch)
			}

			var result *github.RepositoryContentResponse
			var resp *github.Response
			if targetSHA != "" {
				opts.SHA = github.Ptr(targetSHA)
				result, resp, err = client.Repositories.UpdateFile(ctx, owner, repo, repoPath, opts)
			} else {
				result, resp, err = client.Repositories.CreateFile(ctx, owner, repo, repoPath, opts)
			}
			if err != nil {
				return ghErrors.NewGitHubAPIErrorResponse(ctx, "failed to commit sandbox file", resp, err), nil, nil
			}
			defer closeGitHubResponse(resp)
			return marshalWriteResult(map[string]any{
				"result":      result,
				"source_path": filePath,
				"bytes":       len(content),
			})
		},
	)
}

// FileDownload downloads the complete contents of a GitHub repository file into the ChatGPT sandbox.
func FileDownload(t translations.TranslationHelperFunc) inventory.ServerTool {
	return NewTool(
		ToolsetMetadataRepos,
		mcp.Tool{
			Name:        "github_file_download",
			Description: t("TOOL_GITHUB_FILE_DOWNLOAD_DESCRIPTION", "Download the complete bytes of a repository file into the ChatGPT sandbox for local inspection. This is useful when get_file_contents is truncated or unsuitable for very large files."),
			Annotations: &mcp.ToolAnnotations{
				Title:        t("TOOL_GITHUB_FILE_DOWNLOAD_USER_TITLE", "Download repository file"),
				ReadOnlyHint: true,
			},
			InputSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"owner":            {Type: "string", Description: "Repository owner."},
					"repo":             {Type: "string", Description: "Repository name."},
					"path":             {Type: "string", Description: "Repository file path to download."},
					"ref":              {Type: "string", Description: "Optional branch, tag, or commit SHA."},
					"destination_path": {Type: "string", Description: "Optional destination inside /mnt/data. Defaults to /mnt/data/<source basename>."},
					"overwrite":        {Type: "boolean", Description: "Allow replacing an existing sandbox file. Defaults to false."},
				},
				Required: []string{"owner", "repo", "path"},
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
			repoPath, err := RequiredParam[string](args, "path")
			if err != nil {
				return utils.NewToolResultError(err.Error()), nil, nil
			}
			ref, _ := OptionalParam[string](args, "ref")
			overwrite, _ := OptionalParam[bool](args, "overwrite")
			destination, _ := OptionalParam[string](args, "destination_path")
			if destination == "" {
				base := filepath.Base(strings.TrimSpace(repoPath))
				if base == "." || base == "/" || base == "" {
					return utils.NewToolResultError("path must point to a repository file"), nil, nil
				}
				destination = filepath.Join(sandboxRoot(), base)
			}
			destination, err = resolveSandboxDestination(destination)
			if err != nil {
				return utils.NewToolResultError(err.Error()), nil, nil
			}
			if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
				return utils.NewToolResultErrorFromErr("failed to create sandbox destination directory", err), nil, nil
			}

			client, err := deps.GetClient(ctx)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to get GitHub client: %w", err)
			}
			reader, resp, err := client.Repositories.DownloadContents(ctx, owner, repo, repoPath, &github.RepositoryContentGetOptions{Ref: ref})
			if err != nil {
				return utils.NewToolResultErrorFromErr("failed to download repository file", err), nil, nil
			}
			defer reader.Close()
			defer closeGitHubResponse(resp)

			flags := os.O_WRONLY | os.O_CREATE
			if overwrite {
				flags |= os.O_TRUNC
			} else {
				flags |= os.O_EXCL
			}
			out, err := os.OpenFile(destination, flags, 0o644)
			if err != nil {
				if os.IsExist(err) && !overwrite {
					return utils.NewToolResultError("destination_path already exists; set overwrite=true or choose another path"), nil, nil
				}
				return utils.NewToolResultErrorFromErr("failed to open sandbox destination", err), nil, nil
			}
			written, copyErr := io.Copy(out, reader)
			closeErr := out.Close()
			if copyErr != nil {
				_ = os.Remove(destination)
				return utils.NewToolResultErrorFromErr("failed to write downloaded file", copyErr), nil, nil
			}
			if closeErr != nil {
				return utils.NewToolResultErrorFromErr("failed to close downloaded file", closeErr), nil, nil
			}
			return marshalWriteResult(map[string]any{
				"destination_path": destination,
				"bytes":            written,
				"owner":            owner,
				"repo":             repo,
				"path":             repoPath,
				"ref":              ref,
			})
		},
	)
}

func sandboxRoot() string {
	if root := strings.TrimSpace(os.Getenv("GITHUB_MCP_SANDBOX_ROOT")); root != "" {
		return filepath.Clean(root)
	}
	return filepath.Clean("/mnt/data")
}

func pathWithinRoot(root, candidate string) bool {
	rel, err := filepath.Rel(root, candidate)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator))
}

func resolveSandboxPath(path string, requireRegular bool) (string, error) {
	if strings.Contains(path, "://") || !filepath.IsAbs(path) {
		return "", fmt.Errorf("file_path must be an absolute path inside %s", sandboxRoot())
	}
	root, err := filepath.Abs(sandboxRoot())
	if err != nil {
		return "", fmt.Errorf("failed to resolve sandbox root: %w", err)
	}
	candidate, err := filepath.Abs(filepath.Clean(path))
	if err != nil || !pathWithinRoot(root, candidate) {
		return "", fmt.Errorf("file_path must stay inside %s", root)
	}
	resolvedRoot, rootErr := filepath.EvalSymlinks(root)
	resolved, pathErr := filepath.EvalSymlinks(candidate)
	if rootErr == nil && pathErr == nil && !pathWithinRoot(resolvedRoot, resolved) {
		return "", fmt.Errorf("file_path resolves outside %s", root)
	}
	if requireRegular {
		info, err := os.Stat(candidate)
		if err != nil {
			return "", fmt.Errorf("failed to stat file_path: %w", err)
		}
		if !info.Mode().IsRegular() {
			return "", fmt.Errorf("file_path must point to a regular file")
		}
	}
	return candidate, nil
}

func resolveSandboxDestination(path string) (string, error) {
	if strings.Contains(path, "://") || !filepath.IsAbs(path) {
		return "", fmt.Errorf("destination_path must be an absolute path inside %s", sandboxRoot())
	}
	root, err := filepath.Abs(sandboxRoot())
	if err != nil {
		return "", fmt.Errorf("failed to resolve sandbox root: %w", err)
	}
	candidate, err := filepath.Abs(filepath.Clean(path))
	if err != nil || !pathWithinRoot(root, candidate) {
		return "", fmt.Errorf("destination_path must stay inside %s", root)
	}
	parent := filepath.Dir(candidate)
	if resolvedRoot, rootErr := filepath.EvalSymlinks(root); rootErr == nil {
		if resolvedParent, parentErr := filepath.EvalSymlinks(parent); parentErr == nil && !pathWithinRoot(resolvedRoot, resolvedParent) {
			return "", fmt.Errorf("destination_path resolves outside %s", root)
		}
	}
	if info, err := os.Lstat(candidate); err == nil && info.Mode()&os.ModeSymlink != 0 {
		return "", fmt.Errorf("destination_path must not be a symbolic link")
	}
	return candidate, nil
}
