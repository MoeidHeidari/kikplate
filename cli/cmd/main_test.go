package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRootHelpFlag(t *testing.T) {
	rootCmd.SetArgs([]string{"--help"})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("root --help failed: %v", err)
	}
}

func TestConfigCmdHelp(t *testing.T) {
	rootCmd.SetArgs([]string{"config", "--help"})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("config --help failed: %v", err)
	}
}

func TestConfigInitHelp(t *testing.T) {
	rootCmd.SetArgs([]string{"config", "init", "--help"})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("config init --help failed: %v", err)
	}
}

func TestPlatesCmdHelp(t *testing.T) {
	rootCmd.SetArgs([]string{"plates", "--help"})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("plates --help failed: %v", err)
	}
}

func TestPlatesListHelp(t *testing.T) {
	rootCmd.SetArgs([]string{"plates", "list", "--help"})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("plates list --help failed: %v", err)
	}
}

func TestPlatesAddHelp(t *testing.T) {
	rootCmd.SetArgs([]string{"plates", "add", "--help"})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("plates add --help failed: %v", err)
	}
}

func TestLoginCmdHelp(t *testing.T) {
	rootCmd.SetArgs([]string{"login", "--help"})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("login --help failed: %v", err)
	}
}

func TestSearchCmdHelp(t *testing.T) {
	rootCmd.SetArgs([]string{"search", "--help"})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("search --help failed: %v", err)
	}
}

func TestLoadConfigMissing(t *testing.T) {
	_, err := LoadConfig("/nonexistent/path/config.yaml")
	if err == nil {
		t.Error("loading nonexistent config should fail")
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.yaml")

	testConfig := &CLIConfig{
		Server: ServerConfig{
			Address: "https://api.example.com",
		},
		Auth: AuthConfig{
			Token: "test-token-123",
		},
	}

	err := SaveConfig(configPath, testConfig)
	if err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	loaded, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if loaded.Server.Address != testConfig.Server.Address {
		t.Errorf("expected server address %s, got %s", testConfig.Server.Address, loaded.Server.Address)
	}

	if loaded.Auth.Token != testConfig.Auth.Token {
		t.Errorf("expected token %s, got %s", testConfig.Auth.Token, loaded.Auth.Token)
	}
}

func TestLoadConfigMissingServerAddress(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid-config.yaml")

	invalidConfig := &CLIConfig{
		Server: ServerConfig{
			Address: "",
		},
	}

	err := SaveConfig(configPath, invalidConfig)
	if err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	_, err = LoadConfig(configPath)
	if err == nil {
		t.Error("loading config with missing server address should fail")
	}
}

func TestLoadLocalPlatesEmpty(t *testing.T) {
	plates, err := loadLocalPlates()
	if err != nil {
		t.Fatalf("loadLocalPlates failed: %v", err)
	}
	if plates == nil {
		t.Error("expected empty list, got nil")
	}
}

func TestSaveAndLoadLocalPlates(t *testing.T) {
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	testPlates := []LocalPlate{
		{
			Slug:        "test-plate",
			Name:        "Test Plate",
			Description: "A test plate",
			RepoURL:     "https://github.com/test/repo",
			ServerURL:   "https://api.example.com",
		},
	}

	err := saveLocalPlates(testPlates)
	if err != nil {
		t.Fatalf("saveLocalPlates failed: %v", err)
	}

	loaded, err := loadLocalPlates()
	if err != nil {
		t.Fatalf("loadLocalPlates failed: %v", err)
	}

	if len(loaded) != len(testPlates) {
		t.Errorf("expected %d plates, got %d", len(testPlates), len(loaded))
	}

	if loaded[0].Slug != testPlates[0].Slug {
		t.Errorf("expected slug %s, got %s", testPlates[0].Slug, loaded[0].Slug)
	}
}
