package commands_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/FeelsCoderMan/cfgctl/internal/app/cfgctl/commands"
	"github.com/FeelsCoderMan/cfgctl/internal/pkg/storage"
	"github.com/FeelsCoderMan/cfgctl/internal/testutil"
)

func TestNewListCmd_Success(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	cmd := commands.NewListCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(cmd)

	tmpDir := testutil.CreatePackageTempDir(t)
	tmpPath := filepath.Join(tmpDir, "config.json")
	data := map[string]any{"testKey": "testValue"}
	shouldCorruptData := false

	if err := testutil.CreateTmpFile(tmpPath, data, shouldCorruptData); err != nil {
		t.Fatalf("list: failed to create tmp file: %v", err)
	}

	err := testutil.RunCmd(rootCmd, "list", "--path", tmpPath)

	if err != nil {
		t.Fatalf("list: expected no error, got %v", err)
	}

	if mockFileStorage.GetPath() != tmpPath {
		t.Fatalf("list: unexpected path: %s", mockFileStorage.GetPath())
	}

	if !mockFileStorage.IsDataValid(data) {
		t.Fatalf("list: data is not valid")
	}
}

func TestNewListCmd_Failure_MissingPath(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	c := commands.NewListCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(c)

	err := testutil.RunCmd(rootCmd, "list")

	if err == nil {
		t.Fatalf("list: expected error, got nil")
	}

	var cmdError *commands.CommandError
	if !errors.As(err, &cmdError) {
		t.Fatalf("list: expected CommandError, got %T: %v", err, err)
	}

	if cmdError.Kind != commands.KindMissingPath {
		t.Fatalf("list: expected MissingPath error, got %v", cmdError.Kind)
	}
}

func TestNewListCmd_Failure_InvalidData(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	c := commands.NewListCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(c)
	tmpDir := testutil.CreatePackageTempDir(t)

	tmpPath := filepath.Join(tmpDir, "config.json")
	data := map[string]any{"testKey": "testValue"}
	shouldCorruptData := true

	if err := testutil.CreateTmpFile(tmpPath, data, shouldCorruptData); err != nil {
		t.Fatalf("list: failed to create tmp file: %v", err)
	}

	err := testutil.RunCmd(rootCmd, "list", "--path", tmpPath)

	if err == nil {
		t.Fatalf("list: expected error, got nil")
	}

	var cmdError *commands.CommandError
	if !errors.As(err, &cmdError) {
		t.Fatalf("list: expected CommandError, got %T: %v", err, err)
	}

	if cmdError.Kind != commands.KindInvalidJSON {
		t.Fatalf("list: expected InvalidJSON error, got %v", cmdError.Kind)
	}
}

func TestNewListCmd_Failure_MissingFile(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	cmd := commands.NewListCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(cmd)

	tmpDir := testutil.CreatePackageTempDir(t)
	tmpPath := filepath.Join(tmpDir, "config.json")

	err := testutil.RunCmd(rootCmd, "list", "--path", tmpPath)

	if err == nil {
		t.Fatalf("list: expected error on missing file, got nil")
	}

	var cmdError *commands.CommandError
	if !errors.As(err, &cmdError) {
		t.Fatalf("list: expected CommandError, got %v", err)
	}

	if cmdError.Kind != commands.KindFileStoreLoad {
		t.Fatalf("list: expected FileStoreLoad error, got %v", cmdError.Kind)
	}
}

func TestNewListCmd_Failure_InvalidFile(t *testing.T) {
	mockFileStorage := storage.NewMockFileStore()
	cmd := commands.NewListCmd(mockFileStorage.FileStore)
	rootCmd := testutil.NewMockRootCmd()
	rootCmd.AddCommand(cmd)

	tmpDir := testutil.CreatePackageTempDir(t)
	tmpPath := filepath.Join(tmpDir, "config.json")
	data := map[string]any{"name": "Alan", "age": 21}
	shouldCorruptData := true

	if err := testutil.CreateTmpFile(tmpPath, data, shouldCorruptData); err != nil {
		t.Fatalf("list: failed to create tmp file: %v", err)
	}

	err := testutil.RunCmd(rootCmd, "list", "--path", tmpPath)

	if err == nil {
		t.Fatalf("list: expected error, got nil")
	}

	var cmdError *commands.CommandError
	if !errors.As(err, &cmdError) {
		t.Fatalf("list: expected CommandError, got %v", err)
	}

	if cmdError.Kind != commands.KindInvalidJSON {
		t.Fatalf("list: expected InvalidJSON error, got %v", cmdError.Kind)
	}
}
