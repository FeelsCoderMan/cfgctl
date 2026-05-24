package commands_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/FeelsCoderMan/cfgctl/internal/app/cfgctl/commands"
	"github.com/FeelsCoderMan/cfgctl/internal/pkg/storage"
	"github.com/FeelsCoderMan/cfgctl/internal/testutil"
)

func TestNewDeleteCmd_Success(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	cmd := commands.NewDeleteCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(cmd)

	tmpDir := testutil.CreatePackageTempDir(t)
	tmpPath := filepath.Join(tmpDir, "config.json")
	data := map[string]any{"name": "Alan", "age": 21}
	shouldCorruptData := false

	if err := testutil.CreateTmpFile(tmpPath, data, shouldCorruptData); err != nil {
		t.Fatalf("delete: failed to create tmp file: %v", err)
	}

	err := testutil.RunCmd(rootCmd, "delete", "--path", tmpPath, "name")

	if err != nil {
		t.Fatalf("delete: expected no error on execution of delete command, got %v", err)
	}

	value, err := mockFileStorage.FileStore.Get("name")

	if value != nil {
		t.Fatalf("delete: expected nil value, got %v", value)
	}

	if err != nil {
		t.Fatalf("delete: expected no error, got %v", err)
	}
}

func TestNewDeleteCmd_Success_UnpresentKey(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	cmd := commands.NewDeleteCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(cmd)

	tmpDir := testutil.CreatePackageTempDir(t)
	tmpPath := filepath.Join(tmpDir, "config.json")
	data := map[string]any{"name": "Alan", "age": 21}
	shouldCorruptData := false

	if err := testutil.CreateTmpFile(tmpPath, data, shouldCorruptData); err != nil {
		t.Fatalf("delete: failed to create tmp file: %v", err)
	}

	err := testutil.RunCmd(rootCmd, "delete", "--path", tmpPath, "birthday")

	if err != nil {
		t.Fatalf("delete: expected no error on execution of delete command, got %v", err)
	}

	value, err := mockFileStorage.FileStore.Get("birthday")

	if value != nil {
		t.Fatalf("delete: expected nil value, got %v", value)
	}

	if err != nil {
		t.Fatalf("delete: expected no error, got %v", err)
	}
}

func TestNewDeleteCmd_Failure_MissingPath(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	cmd := commands.NewDeleteCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(cmd)

	tmpDir := testutil.CreatePackageTempDir(t)
	tmpPath := filepath.Join(tmpDir, "config.json")
	data := map[string]any{"name": "Alan", "age": 21}
	shouldCorruptData := false

	if err := testutil.CreateTmpFile(tmpPath, data, shouldCorruptData); err != nil {
		t.Fatalf("delete: failed to create tmp file: %v", err)
	}

	err := testutil.RunCmd(rootCmd, "delete", "name")

	if err == nil {
		t.Fatalf("delete: expected error on missing path, got nil")
	}

	var cmdError *commands.CommandError
	if !errors.As(err, &cmdError) {
		t.Fatalf("delete: expected CommandError, got %v", err)
	}

	if cmdError.Kind != commands.KindMissingPath {
		t.Fatalf("delete: expected MissingPath error, got %v", cmdError.Kind)
	}
}

func TestNewDeleteCmd_Failure_MissingKey(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	cmd := commands.NewDeleteCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(cmd)

	tmpDir := testutil.CreatePackageTempDir(t)
	tmpPath := filepath.Join(tmpDir, "config.json")
	data := map[string]any{"name": "Alan", "age": 21}
	shouldCorruptData := false

	if err := testutil.CreateTmpFile(tmpPath, data, shouldCorruptData); err != nil {
		t.Fatalf("delete: failed to create tmp file: %v", err)
	}

	err := testutil.RunCmd(rootCmd, "delete", "--path", tmpPath)

	if err == nil {
		t.Fatalf("delete: expected error on missing key, got nil")
	}

	var cmdError *commands.CommandError
	if !errors.As(err, &cmdError) {
		t.Fatalf("delete: expected CommandError, got %v", err)
	}

	if cmdError.Kind != commands.KindMissingArgs {
		t.Fatalf("delete: expected MissingArgs error, got %v", cmdError.Kind)
	}
}

func TestNewDeleteCmd_Failure_MissingFile(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	cmd := commands.NewDeleteCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(cmd)

	tmpDir := testutil.CreatePackageTempDir(t)
	tmpPath := filepath.Join(tmpDir, "config.json")

	err := testutil.RunCmd(rootCmd, "delete", "--path", tmpPath, "name")

	if err == nil {
		t.Fatalf("delete: expected error on missing file, got nil")
	}

	var cmdError *commands.CommandError
	if !errors.As(err, &cmdError) {
		t.Fatalf("delete: expected CommandError, got %v", err)
	}

	if cmdError.Kind != commands.KindFileStoreLoad {
		t.Fatalf("delete: expected FileStoreLoad error, got %v", cmdError.Kind)
	}
}

func TestNewDeleteCmd_Failure_InvalidFile(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	cmd := commands.NewDeleteCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(cmd)

	tmpDir := testutil.CreatePackageTempDir(t)
	tmpPath := filepath.Join(tmpDir, "config.json")
	data := map[string]any{"name": "Alan", "age": 21}
	shouldCorruptData := true

	if err := testutil.CreateTmpFile(tmpPath, data, shouldCorruptData); err != nil {
		t.Fatalf("delete: failed to create tmp file: %v", err)
	}

	err := testutil.RunCmd(rootCmd, "delete", "--path", tmpPath, "name")

	if err == nil {
		t.Fatalf("delete: expected error on missing key, got nil")
	}

	var cmdError *commands.CommandError
	if !errors.As(err, &cmdError) {
		t.Fatalf("delete: expected CommandError, got %v", err)
	}

	if cmdError.Kind != commands.KindInvalidJSON {
		t.Fatalf("delete: expected InvalidJSON error, got %v", cmdError.Kind)
	}
}
