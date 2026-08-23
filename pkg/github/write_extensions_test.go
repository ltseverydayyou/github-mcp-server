package github

import (
	"net/url"
	"testing"
)

func TestValidateRelativeGitHubEndpoint(t *testing.T) {
	t.Parallel()

	valid, err := validateRelativeGitHubEndpoint("/repos/octo/example/releases/1/assets?name=a.zip")
	if err != nil {
		t.Fatalf("validateRelativeGitHubEndpoint returned error: %v", err)
	}
	if valid != "repos/octo/example/releases/1/assets?name=a.zip" {
		t.Fatalf("unexpected endpoint: %q", valid)
	}

	for _, endpoint := range []string{
		"https://uploads.github.com/repos/octo/example/releases/1/assets",
		"../repos/octo/example",
		"",
	} {
		if _, err := validateRelativeGitHubEndpoint(endpoint); err == nil {
			t.Errorf("expected endpoint %q to be rejected", endpoint)
		}
	}
}

func TestAddUploadQuery(t *testing.T) {
	t.Parallel()

	endpoint, err := addUploadQuery(
		"repos/octo/example/releases/1/assets?existing=1",
		"asset.zip",
		"binary",
		map[string]any{"foo": "bar", "multi": []any{"a", "b"}},
	)
	if err != nil {
		t.Fatalf("addUploadQuery returned error: %v", err)
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		t.Fatalf("failed to parse result: %v", err)
	}
	q := u.Query()
	if q.Get("name") != "asset.zip" || q.Get("label") != "binary" || q.Get("existing") != "1" || q.Get("foo") != "bar" {
		t.Fatalf("unexpected query parameters: %v", q)
	}
	multi := q["multi"]
	if len(multi) != 2 || multi[0] != "a" || multi[1] != "b" {
		t.Fatalf("unexpected multi query parameter: %v", multi)
	}
}
