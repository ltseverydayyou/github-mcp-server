package github

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	ghErrors "github.com/github/github-mcp-server/pkg/errors"
	"github.com/github/github-mcp-server/pkg/inventory"
	"github.com/github/github-mcp-server/pkg/scopes"
	"github.com/github/github-mcp-server/pkg/translations"
	"github.com/github/github-mcp-server/pkg/utils"
	"github.com/google/go-github/v89/github"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const maxChatGPTRepositoryFileBytes int64 = 100 << 20

func generateReposToolsetInstructions(_ *inventory.Inventory) string {
	return `## Repository files

When a ChatGPT attachment, generated artifact, or sandbox file must be created or replaced in a repository, prefer github_repository_file_upload. Pass the local sandbox file to its file parameter; ChatGPT supplies the actual file through openai/fileParams. Do not create a GitHub Actions workflow merely to transport or update file bytes when this tool is available.

When complete repository file bytes are needed for local inspection, especially for a large file, prefer github_file_download. It returns an MCP resource link that lets the host materialize the whole file without placing its contents in model context.`
}

type chatGPTFileInput struct {
	DownloadURL string
	FileID      string
	MIMEType    string
	FileName    string
}

func chatGPTFileSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "object",
		Description: "ChatGPT-provided file object. Pass a local sandbox file to this parameter; ChatGPT supplies the authorized file reference.",
		Properties: map[string]*jsonschema.Schema{
			"download_url": {Type: "string", Description: "Temporary authorized HTTPS download URL supplied by ChatGPT."},
			"file_id":      {Type: "string", Description: "ChatGPT file identifier."},
			"mime_type":    {Type: "string", Description: "Optional MIME type supplied by ChatGPT."},
			"file_name":    {Type: "string", Description: "Optional original file name supplied by ChatGPT."},
		},
		Required: []string{"download_url", "file_id"},
	}
}

func parseChatGPTFileInput(args map[string]any, name string) (chatGPTFileInput, error) {
	raw, ok := args[name]
	if !ok || raw == nil {
		return chatGPTFileInput{}, fmt.Errorf("missing required parameter: %s", name)
	}
	obj, ok := raw.(map[string]any)
	if !ok {
		return chatGPTFileInput{}, fmt.Errorf("%s must be the ChatGPT-provided file object", name)
	}
	stringField := func(key string) string {
		value, _ := obj[key].(string)
		return strings.TrimSpace(value)
	}
	file := chatGPTFileInput{
		DownloadURL: stringField("download_url"),
		FileID:      stringField("file_id"),
		MIMEType:    stringField("mime_type"),
		FileName:    stringField("file_name"),
	}
	if file.DownloadURL == "" {
		return chatGPTFileInput{}, fmt.Errorf("%s.download_url is required", name)
	}
	if file.FileID == "" {
		return chatGPTFileInput{}, fmt.Errorf("%s.file_id is required", name)
	}
	return file, nil
}

func unsafeDownloadIP(ip net.IP) bool {
	return ip == nil || ip.IsUnspecified() || ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsMulticast()
}

func validateChatGPTDownloadURL(raw string) (*url.URL, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("invalid ChatGPT file download URL: %w", err)
	}
	if u.Scheme != "https" || u.Hostname() == "" {
		return nil, fmt.Errorf("ChatGPT file download URL must use HTTPS")
	}
	if u.User != nil {
		return nil, fmt.Errorf("ChatGPT file download URL must not contain embedded credentials")
	}
	if ip := net.ParseIP(u.Hostname()); ip != nil && unsafeDownloadIP(ip) {
		return nil, fmt.Errorf("ChatGPT file download URL resolves to a non-public address")
	}
	return u, nil
}

func chatGPTFileHTTPClient() *http.Client {
	dialer := &net.Dialer{Timeout: 30 * time.Second, KeepAlive: 30 * time.Second}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.Proxy = nil
	transport.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
		host, port, err := net.SplitHostPort(address)
		if err != nil {
			return nil, err
		}
		ips, err := net.DefaultResolver.LookupIP(ctx, "ip", host)
		if err != nil {
			return nil, err
		}
		if len(ips) == 0 {
			return nil, fmt.Errorf("download host did not resolve")
		}
		for _, ip := range ips {
			if unsafeDownloadIP(ip) {
				return nil, fmt.Errorf("download host resolves to a non-public address")
			}
		}
		var lastErr error
		for _, ip := range ips {
			conn, dialErr := dialer.DialContext(ctx, network, net.JoinHostPort(ip.String(), port))
			if dialErr == nil {
				return conn, nil
			}
			lastErr = dialErr
		}
		return nil, lastErr
	}
	return &http.Client{
		Transport: transport,
		Timeout:   2 * time.Minute,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 5 {
				return fmt.Errorf("too many redirects while downloading ChatGPT file")
			}
			_, err := validateChatGPTDownloadURL(req.URL.String())
			return err
		},
	}
}

func downloadChatGPTFile(ctx context.Context, rawURL string) ([]byte, error) {
	u, err := validateChatGPTDownloadURL(rawURL)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create ChatGPT file request: %w", err)
	}
	resp, err := chatGPTFileHTTPClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download ChatGPT file: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("ChatGPT file download returned HTTP %d", resp.StatusCode)
	}
	if resp.ContentLength > maxChatGPTRepositoryFileBytes {
		return nil, fmt.Errorf("ChatGPT file is too large: maximum supported size is %d bytes", maxChatGPTRepositoryFileBytes)
	}
	content, err := io.ReadAll(io.LimitReader(resp.Body, maxChatGPTRepositoryFileBytes+1))
	if err != nil {
		return nil, fmt.Errorf("failed to read ChatGPT file: %w", err)
	}
	if int64(len(content)) > maxChatGPTRepositoryFileBytes {
		return nil, fmt.Errorf("ChatGPT file is too large: maximum supported size is %d bytes", maxChatGPTRepositoryFileBytes)
	}
	return content, nil
}

// RepositoryFileUpload commits a ChatGPT-provided file directly into a repository.
func RepositoryFileUpload(t translations.TranslationHelperFunc) inventory.ServerTool {
	return NewTool(
		ToolsetMetadataRepos,
		mcp.Tool{
			Meta: mcp.Meta{"openai/fileParams": []string{"file"}},
			Name: "github_repository_file_upload",
			Description: t("TOOL_GITHUB_REPOSITORY_FILE_UPLOAD_DESCRIPTION", "Create or replace a repository file from a ChatGPT attachment, generated artifact, or sandbox file without sending the file contents through model context. ChatGPT injects the file parameter through openai/fileParams; pass the local /mnt/data file to that parameter rather than a pathname string. The tool automatically resolves the existing target SHA when updating."),
			Annotations: &mcp.ToolAnnotations{
				Title:           t("TOOL_GITHUB_REPOSITORY_FILE_UPLOAD_USER_TITLE", "Commit ChatGPT file to repository"),
				ReadOnlyHint:    false,
				DestructiveHint: jsonschema.Ptr(true),
			},
			InputSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"owner":   {Type: "string", Description: "Repository owner."},
					"repo":    {Type: "string", Description: "Repository name."},
					"path":    {Type: "string", Description: "Repository path to create or replace."},
					"file":    chatGPTFileSchema(),
					"message": {Type: "string", Description: "Commit message."},
					"branch":  {Type: "string", Description: "Branch to update. Omit to use the repository default branch."},
					"sha":     {Type: "string", Description: "Optional expected blob SHA for an existing target. When omitted, the current SHA is detected automatically."},
				},
				Required: []string{"owner", "repo", "path", "file", "message"},
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
			message, err := RequiredParam[string](args, "message")
			if err != nil {
				return utils.NewToolResultError(err.Error()), nil, nil
			}
			file, err := parseChatGPTFileInput(args, "file")
			if err != nil {
				return utils.NewToolResultError(err.Error()), nil, nil
			}
			if strings.TrimSpace(repoPath) == "" || strings.HasSuffix(repoPath, "/") {
				return utils.NewToolResultError("path must point to a repository file"), nil, nil
			}
			if strings.TrimSpace(message) == "" {
				return utils.NewToolResultError("message must not be empty"), nil, nil
			}

			content, err := downloadChatGPTFile(ctx, file.DownloadURL)
			if err != nil {
				return utils.NewToolResultError(err.Error()), nil, nil
			}

			branch, _ := OptionalParam[string](args, "branch")
			targetSHA, _ := OptionalParam[string](args, "sha")
			client, err := deps.GetClient(ctx)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to get GitHub client: %w", err)
			}

			if targetSHA == "" {
				existing, directory, resp, getErr := client.Repositories.GetContents(ctx, owner, repo, repoPath, &github.RepositoryContentGetOptions{Ref: branch})
				switch {
				case getErr == nil:
					defer closeGitHubResponse(resp)
					if existing == nil || directory != nil {
						return utils.NewToolResultError("path resolves to a directory, not a file"), nil, nil
					}
					targetSHA = existing.GetSHA()
				case resp == nil || resp.StatusCode != http.StatusNotFound:
					return ghErrors.NewGitHubAPIErrorResponse(ctx, "failed to inspect repository target", resp, getErr), nil, nil
				default:
					closeGitHubResponse(resp)
				}
			}

			opts := &github.RepositoryContentFileOptions{Message: github.Ptr(message), Content: content}
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
				return ghErrors.NewGitHubAPIErrorResponse(ctx, "failed to commit ChatGPT file", resp, err), nil, nil
			}
			defer closeGitHubResponse(resp)
			return marshalWriteResult(map[string]any{
				"result":    result,
				"file_id":   file.FileID,
				"file_name": file.FileName,
				"mime_type": file.MIMEType,
				"bytes":     len(content),
			})
		},
	)
}

// FileDownload returns a complete repository file as an MCP resource link for host-side materialization.
func FileDownload(t translations.TranslationHelperFunc) inventory.ServerTool {
	return NewTool(
		ToolsetMetadataRepos,
		mcp.Tool{
			Name:        "github_file_download",
			Description: t("TOOL_GITHUB_FILE_DOWNLOAD_DESCRIPTION", "Return a complete repository file as an MCP resource link so the host can materialize or download the full bytes without putting large source content in model context. Prefer this when get_file_contents is too large or when the whole file is needed for local inspection."),
			Annotations: &mcp.ToolAnnotations{
				Title:        t("TOOL_GITHUB_FILE_DOWNLOAD_USER_TITLE", "Download complete repository file"),
				ReadOnlyHint: true,
			},
			InputSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"owner": {Type: "string", Description: "Repository owner."},
					"repo":  {Type: "string", Description: "Repository name."},
					"path":  {Type: "string", Description: "Repository file path to download."},
					"ref":   {Type: "string", Description: "Optional git ref such as a branch, tag, or refs/pull/<number>/head."},
					"sha":   {Type: "string", Description: "Optional commit SHA. Takes precedence over ref."},
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
			sha, _ := OptionalParam[string](args, "sha")
			if strings.TrimSpace(repoPath) == "" || strings.HasSuffix(repoPath, "/") {
				return utils.NewToolResultError("path must point to a repository file"), nil, nil
			}

			client, err := deps.GetClient(ctx)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to get GitHub client: %w", err)
			}
			rawOpts, _, err := resolveGitReference(ctx, client, owner, repo, ref, sha)
			if err != nil {
				return utils.NewToolResultError(fmt.Sprintf("failed to resolve git reference: %s", err)), nil, nil
			}
			resolvedRef := ref
			if rawOpts.SHA != "" {
				resolvedRef = rawOpts.SHA
			}
			fileContent, dirContent, resp, err := client.Repositories.GetContents(ctx, owner, repo, repoPath, &github.RepositoryContentGetOptions{Ref: resolvedRef})
			if err != nil {
				return ghErrors.NewGitHubAPIErrorResponse(ctx, "failed to inspect repository file", resp, err), nil, nil
			}
			defer closeGitHubResponse(resp)
			if fileContent == nil || dirContent != nil {
				return utils.NewToolResultError("path resolves to a directory, not a file"), nil, nil
			}

			resourceURI, err := expandRepoResourceURI(owner, repo, sha, resolvedRef, strings.Split(strings.Trim(repoPath, "/"), "/"))
			if err != nil {
				return utils.NewToolResultError("failed to build repository resource URI"), nil, nil
			}
			size := int64(fileContent.GetSize())
			resource := &mcp.ResourceLink{
				URI:   resourceURI,
				Name:  fileContent.GetName(),
				Title: fmt.Sprintf("File: %s", repoPath),
				Size:  &size,
			}
			message := fmt.Sprintf("Complete repository file available as a downloadable MCP resource (%d bytes, blob %s).", size, fileContent.GetSHA())
			return utils.NewToolResultResourceLink(message, resource), nil, nil
		},
	)
}
