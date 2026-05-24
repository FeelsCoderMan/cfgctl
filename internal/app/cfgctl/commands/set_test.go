package commands_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/FeelsCoderMan/cfgctl/internal/app/cfgctl/commands"
	"github.com/FeelsCoderMan/cfgctl/internal/pkg/storage"
	"github.com/FeelsCoderMan/cfgctl/internal/testutil"
)

func TestNewSetCmd_Success(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	cmd := commands.NewSetCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(cmd)

	tmpDir := testutil.CreatePackageTempDir(t)
	tmpPath := filepath.Join(tmpDir, "config.json")
	data := map[string]any{"name": "Alan", "age": 21}
	shouldCorruptData := false

	if err := testutil.CreateTmpFile(tmpPath, data, shouldCorruptData); err != nil {
		t.Fatalf("set: failed to create tmp file: %v", err)
	}

	key := "name"
	err := testutil.RunCmd(rootCmd, "set", "--path", tmpPath, key, "Bob")

	if err != nil {
		t.Fatalf("set: expected no error on execution of set command, got %v", err)
	}

	value, err := mockFileStorage.FileStore.Get(key)

	if err != nil {
		t.Fatalf("set: expected no error, got %v", err)
	}

	if value != "Bob" {
		t.Fatalf("set: expected value 'Bob', got %v", value)
	}
}

func TestNewSetCmd_Failure_MissingPath(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	cmd := commands.NewSetCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(cmd)

	tmpDir := testutil.CreatePackageTempDir(t)
	tmpPath := filepath.Join(tmpDir, "config.json")
	data := map[string]any{"name": "Alan", "age": 21}
	shouldCorruptData := false

	if err := testutil.CreateTmpFile(tmpPath, data, shouldCorruptData); err != nil {
		t.Fatalf("set: failed to create tmp file: %v", err)
	}

	err := testutil.RunCmd(rootCmd, "set", "name", "Bob")

	if err == nil {
		t.Fatalf("set: expected error on missing path, got nil")
	}

	var cmdError *commands.CommandError
	if !errors.As(err, &cmdError) {
		t.Fatalf("set: expected CommandError, got %v", err)
	}

	if cmdError.Kind != commands.KindMissingPath {
		t.Fatalf("set: expected MissingPath error, got %v", cmdError.Kind)
	}
}

func TestNewSetCmd_Failure_MissingKey(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	cmd := commands.NewSetCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(cmd)

	tmpDir := testutil.CreatePackageTempDir(t)
	tmpPath := filepath.Join(tmpDir, "config.json")
	data := map[string]any{"name": "Alan", "age": 21}
	shouldCorruptData := false

	if err := testutil.CreateTmpFile(tmpPath, data, shouldCorruptData); err != nil {
		t.Fatalf("set: failed to create tmp file: %v", err)
	}

	err := testutil.RunCmd(rootCmd, "set", "--path", tmpPath, "Bob")

	if err == nil {
		t.Fatalf("set: expected error on missing key, got nil")
	}

	var cmdError *commands.CommandError
	if !errors.As(err, &cmdError) {
		t.Fatalf("set: expected CommandError, got %v", err)
	}

	if cmdError.Kind != commands.KindMissingArgs {
		t.Fatalf("set: expected MissingArgs error, got %v", cmdError.Kind)
	}
}

func TestNewSetCmd_Failure_MissingFile(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	cmd := commands.NewSetCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(cmd)

	tmpDir := testutil.CreatePackageTempDir(t)
	tmpPath := filepath.Join(tmpDir, "config.json")

	err := testutil.RunCmd(rootCmd, "set", "--path", tmpPath, "name", "Bob")

	if err == nil {
		t.Fatalf("set: expected error on missing file, got nil")
	}

	var cmdError *commands.CommandError
	if !errors.As(err, &cmdError) {
		t.Fatalf("set: expected CommandError, got %v", err)
	}

	if cmdError.Kind != commands.KindFileStoreLoad {
		t.Fatalf("set: expected FileStoreLoad error, got %v", cmdError.Kind)
	}
}

func TestNewSetCmd_Failure_InvalidFile(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	cmd := commands.NewSetCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(cmd)

	tmpDir := testutil.CreatePackageTempDir(t)
	tmpPath := filepath.Join(tmpDir, "config.json")
	data := map[string]any{"name": "Alan", "age": 21}
	shouldCorruptData := true

	if err := testutil.CreateTmpFile(tmpPath, data, shouldCorruptData); err != nil {
		t.Fatalf("set: failed to create tmp file: %v", err)
	}

	err := testutil.RunCmd(rootCmd, "set", "--path", tmpPath, "name", "Bob")

	if err == nil {
		t.Fatalf("set: expected error on invalid file, got nil")
	}

	var cmdError *commands.CommandError
	if !errors.As(err, &cmdError) {
		t.Fatalf("set: expected CommandError, got %v", err)
	}

	if cmdError.Kind != commands.KindInvalidJSON {
		t.Fatalf("set: expected InvalidJSON error, got %v", cmdError.Kind)
	}
}
