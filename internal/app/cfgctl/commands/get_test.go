package commands_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/FeelsCoderMan/cfgctl/internal/app/cfgctl/commands"
	"github.com/FeelsCoderMan/cfgctl/internal/pkg/storage"
	"github.com/FeelsCoderMan/cfgctl/internal/testutil"
)

func TestNewGetCmd_Success(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	cmd := commands.NewGetCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(cmd)

	tmpDir := testutil.CreatePackageTempDir(t)
	tmpPath := filepath.Join(tmpDir, "config.json")
	data := map[string]any{"name": "Alan", "age": 21}
	shouldCorruptData := false

	if err := testutil.CreateTmpFile(tmpPath, data, shouldCorruptData); err != nil {
		t.Fatalf("get: failed to create tmp file: %v", err)
	}

	key := "name"
	err := testutil.RunCmd(rootCmd, "get", "--path", tmpPath, key)

	if err != nil {
		t.Fatalf("get: expected no error on execution of get command, got %v", err)
	}

	value, err := mockFileStorage.FileStore.Get(key)

	if err != nil {
		t.Fatalf("get: expected no error, got %v", err)
	}

	if value != "Alan" {
		t.Fatalf("get: expected value 'Alan', got %v", value)
	}

}

func TestNewGetCmd_Success_UnpresentKey(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	cmd := commands.NewGetCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(cmd)

	tmpDir := testutil.CreatePackageTempDir(t)
	tmpPath := filepath.Join(tmpDir, "config.json")
	data := map[string]any{"name": "Alan", "age": 21}
	shouldCorruptData := false

	if err := testutil.CreateTmpFile(tmpPath, data, shouldCorruptData); err != nil {
		t.Fatalf("get: failed to create tmp file: %v", err)
	}

	key := "birthday"
	err := testutil.RunCmd(rootCmd, "get", "--path", tmpPath, key)

	if err != nil {
		t.Fatalf("get: expected no error on execution of get command, got %v", err)
	}

	value, err := mockFileStorage.FileStore.Get(key)

	if err != nil {
		t.Fatalf("get: expected no error, got %v", err)
	}

	if value != nil {
		t.Fatalf("get: expected nil value, got %v", value)
	}
}

func TestNewGetCmd_Failure_MissingPath(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	cmd := commands.NewGetCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(cmd)

	tmpDir := testutil.CreatePackageTempDir(t)
	tmpPath := filepath.Join(tmpDir, "config.json")
	data := map[string]any{"name": "Alan", "age": 21}
	shouldCorruptData := false

	if err := testutil.CreateTmpFile(tmpPath, data, shouldCorruptData); err != nil {
		t.Fatalf("get: failed to create tmp file: %v", err)
	}

	err := testutil.RunCmd(rootCmd, "get", "name")

	if err == nil {
		t.Fatalf("get: expected error on missing path, got nil")
	}

	var cmdError *commands.CommandError
	if !errors.As(err, &cmdError) {
		t.Fatalf("get: expected CommandError, got %v", err)
	}

	if cmdError.Kind != commands.KindMissingPath {
		t.Fatalf("get: expected MissingPath error, got %v", cmdError.Kind)
	}
}

func TestNewGetCmd_Failure_MissingKey(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	cmd := commands.NewGetCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(cmd)

	tmpDir := testutil.CreatePackageTempDir(t)
	tmpPath := filepath.Join(tmpDir, "config.json")
	data := map[string]any{"name": "Alan", "age": 21}
	shouldCorruptData := false

	if err := testutil.CreateTmpFile(tmpPath, data, shouldCorruptData); err != nil {
		t.Fatalf("get: failed to create tmp file: %v", err)
	}

	err := testutil.RunCmd(rootCmd, "get", "--path", tmpPath)

	if err == nil {
		t.Fatalf("get: expected error on missing key, got nil")
	}

	var cmdError *commands.CommandError
	if !errors.As(err, &cmdError) {
		t.Fatalf("get: expected CommandError, got %v", err)
	}

	if cmdError.Kind != commands.KindMissingArgs {
		t.Fatalf("get: expected MissingArgs error, got %v", cmdError.Kind)
	}
}

func TestNewGetCmd_Failure_MissingFile(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	cmd := commands.NewGetCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(cmd)

	tmpDir := testutil.CreatePackageTempDir(t)
	tmpPath := filepath.Join(tmpDir, "config.json")

	err := testutil.RunCmd(rootCmd, "get", "--path", tmpPath, "name")

	if err == nil {
		t.Fatalf("get: expected error on missing file, got nil")
	}

	var cmdError *commands.CommandError
	if !errors.As(err, &cmdError) {
		t.Fatalf("get: expected CommandError, got %v", err)
	}

	if cmdError.Kind != commands.KindFileStoreLoad {
		t.Fatalf("get: expected FileStoreLoad error, got %v", cmdError.Kind)
	}
}

func TestNewGetCmd_Failure_InvalidFile(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	cmd := commands.NewGetCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(cmd)

	tmpDir := testutil.CreatePackageTempDir(t)
	tmpPath := filepath.Join(tmpDir, "config.json")
	data := map[string]any{"name": "Alan", "age": 21}
	shouldCorruptData := true

	if err := testutil.CreateTmpFile(tmpPath, data, shouldCorruptData); err != nil {
		t.Fatalf("get: failed to create tmp file: %v", err)
	}

	err := testutil.RunCmd(rootCmd, "get", "--path", tmpPath, "name")

	if err == nil {
		t.Fatalf("get: expected error on missing key, got nil")
	}

	var cmdError *commands.CommandError
	if !errors.As(err, &cmdError) {
		t.Fatalf("get: expected CommandError, got %v", err)
	}

	if cmdError.Kind != commands.KindInvalidJSON {
		t.Fatalf("get: expected InvalidJSON error, got %v", cmdError.Kind)
	}
}
