package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestConfigDir(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	// darwin branch
	if runtime.GOOS == "darwin" {
		dir, err := configDir()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := filepath.Join(tmpDir, "Library", "Application Support", "slackernews")
		if dir != expected {
			t.Errorf("expected %q, got %q", expected, dir)
		}
	}

	// linux/default branch
	if runtime.GOOS == "linux" {
		dir, err := configDir()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := filepath.Join(tmpDir, ".config", "slackernews")
		if dir != expected {
			t.Errorf("expected %q, got %q", expected, dir)
		}
	}
}

func TestConfigDirWindows(t *testing.T) {
	if runtime.GOOS == "windows" {
		tmpDir := t.TempDir()
		t.Setenv("APPDATA", tmpDir)

		dir, err := configDir()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := filepath.Join(tmpDir, "slackernews")
		if dir != expected {
			t.Errorf("expected %q, got %q", expected, dir)
		}
	}
}

func TestConfigDirWindowsMissingAppData(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Setenv("APPDATA", "")

		_, err := configDir()
		if err == nil {
			t.Fatal("expected error when APPDATA is not set")
		}
	}
}

func TestLoadEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.InstanceURL != "" {
		t.Errorf("expected empty InstanceURL, got %q", cfg.InstanceURL)
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	cfg := Config{InstanceURL: "https://slackernews.example.com"}
	if err := Save(cfg); err != nil {
		t.Fatalf("unexpected error saving: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("unexpected error loading: %v", err)
	}
	if loaded.InstanceURL != cfg.InstanceURL {
		t.Errorf("expected %q, got %q", cfg.InstanceURL, loaded.InstanceURL)
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	// Write invalid JSON to config file
	dir, _ := configDir()
	os.MkdirAll(dir, 0755)
	path := filepath.Join(dir, "config.json")
	os.WriteFile(path, []byte("{not json"), 0600)

	_, err := Load()
	if err == nil {
		t.Fatal("expected error loading invalid JSON")
	}
}

func TestLoadValidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	dir, _ := configDir()
	os.MkdirAll(dir, 0755)
	path := filepath.Join(dir, "config.json")
	data, _ := json.Marshal(Config{InstanceURL: "https://test.com"})
	os.WriteFile(path, data, 0600)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.InstanceURL != "https://test.com" {
		t.Errorf("expected %q, got %q", "https://test.com", cfg.InstanceURL)
	}
}
