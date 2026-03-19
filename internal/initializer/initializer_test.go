package initializer

import (
	"os"
	"path/filepath"
	"testing"
)

func setupInitializer(t *testing.T) string {
	t.Helper()

	original, _ := os.Getwd()
	tmp := t.TempDir()
	
	os.Chdir(tmp)
	t.Cleanup(func() {
		os.Chdir(original)
	})
	
	return tmp
}

func TestInitialize_Success(t *testing.T) {
	tmp := setupInitializer(t)

	projectName := "test-project"
	err := Initialize(projectName)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []string{
		"content",
		"content/blogs",
		"layout",
		"layout/default",
		"layout/partials",
		"config.yaml",
	}

	for _, name := range expected {
		path := filepath.Join(tmp, projectName, name)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("expected %s to exist", name)
		}
	}
}

func TestInitialize_FailsIfDirExists(t *testing.T) {
	_ = setupInitializer(t)

	projectName := "test-project"
	if err := os.Mkdir(projectName, 0755); err != nil {
		t.Fatalf("failed to create dir: %v", err)
	}

	err := Initialize(projectName)
	if err == nil {
		t.Fatal("expected to get error when directory already exists, got nil")
	}
}