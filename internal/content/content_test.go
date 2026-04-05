package content

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestContentCreate_Success(t *testing.T) {
	tmp := t.TempDir()
	path := "/blogs/test-blog.md" 

	outPath, err := Create(path, tmp)
	if err != nil {
		t.Fatalf("failed to create content: %v", err)
	}

	expectedOutPath := filepath.Join(tmp, path)

	if outPath != expectedOutPath {
		t.Errorf("expected outPath to be %s, got %s", expectedOutPath, outPath)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("failed to read file %s: %v", outPath, err)
	}

	content := string(data)

	if !strings.Contains(content, "---") {
		t.Error("expected front matter delimiter")
	}
	if !strings.Contains(content, "title: test blog") {
		t.Error("expected title in front matter")
	}
	if !strings.Contains(content, "date:") {
		t.Error("expected date in front matter")
	}
	if !strings.Contains(content, "# test blog") {
		t.Error("expected title in front matter")
	}
}
func TestContentCreate_NoExtension(t *testing.T) {
	tmp := t.TempDir()
	path := "/blogs/another-test-blog" 

	outPath, err := Create(path, tmp)
	if err != nil {
		t.Fatalf("failed to create content: %v", err)
	}

	if filepath.Ext(outPath) != ".md" {
		t.Errorf("expectd .md extension, got %s", filepath.Ext(outPath))
	}
}

func TestContentCreate_InvalidExtension(t *testing.T) {
	tmp := t.TempDir()
	path := "/blogs/invalid-blog.txt"

	outPath, err := Create(path, tmp)
	if err != nil {
		t.Fatalf("faild to create content: %v", err)
	}

	if filepath.Ext(outPath) != ".md" {
		t.Errorf("expected extension to be .md, got %s", filepath.Ext(outPath))
	}
}

func TestContentCreate_MakeNestedDir(t *testing.T) {
	tmp := t.TempDir()
	path := "/deep/nested/dir/blog.md"

	_, err := Create(path, tmp)
	if err != nil {
		t.Fatalf("faild to create content: %v", err)
	}

	dirPath := filepath.Join(tmp, "deep", "nested", "dir")
	info, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			t.Errorf("expected directory to exist: %v", err)
		} else {
			t.Fatalf("failed to stat %s: %v", dirPath, err)
		}
	}

	if !info.IsDir() {
		t.Errorf("expected path to be a directory")
	}
}

func TestNormalizePath(t *testing.T) {
	tests := []struct{
		input string
		expected string
	}{
		{"post.md", "post.md"},
		{"post.txt", "post.md"},
		{"post", "post.md"},
		{"/blogs/post.md", "/blogs/post.md"},
		{"/blogs/post", "/blogs/post.md"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := normalizePath(tt.input)
			if result != tt.expected {
				t.Errorf("expected path to be %s, got %s", tt.expected, result)
			}
		})
	}
}