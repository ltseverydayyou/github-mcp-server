package github

import (
	"net"
	"testing"
)

func TestParseChatGPTFileInput(t *testing.T) {
	file, err := parseChatGPTFileInput(map[string]any{
		"file": map[string]any{
			"download_url": "https://example.com/file",
			"file_id":      "file_123",
			"mime_type":    "text/plain",
			"file_name":    "large.lua",
		},
	}, "file")
	if err != nil {
		t.Fatalf("parseChatGPTFileInput returned error: %v", err)
	}
	if file.DownloadURL != "https://example.com/file" || file.FileID != "file_123" || file.FileName != "large.lua" {
		t.Fatalf("unexpected parsed file: %#v", file)
	}
}

func TestParseChatGPTFileInputRequiresProvidedFileFields(t *testing.T) {
	_, err := parseChatGPTFileInput(map[string]any{
		"file": map[string]any{"file_id": "file_123"},
	}, "file")
	if err == nil {
		t.Fatal("expected missing download_url to be rejected")
	}
}

func TestValidateChatGPTDownloadURL(t *testing.T) {
	if _, err := validateChatGPTDownloadURL("https://example.com/file"); err != nil {
		t.Fatalf("expected public HTTPS URL to be accepted: %v", err)
	}
	for _, raw := range []string{
		"http://example.com/file",
		"https://user:pass@example.com/file",
		"https://127.0.0.1/file",
		"https://10.0.0.1/file",
	} {
		if _, err := validateChatGPTDownloadURL(raw); err == nil {
			t.Fatalf("expected URL to be rejected: %s", raw)
		}
	}
}

func TestUnsafeDownloadIP(t *testing.T) {
	if !unsafeDownloadIP(net.ParseIP("127.0.0.1")) {
		t.Fatal("loopback must be rejected")
	}
	if !unsafeDownloadIP(net.ParseIP("192.168.1.1")) {
		t.Fatal("private address must be rejected")
	}
	if unsafeDownloadIP(net.ParseIP("8.8.8.8")) {
		t.Fatal("public address must be accepted")
	}
}
