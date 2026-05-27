package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Logger struct {
	log *log.Logger
}

func findProjectRoot(start string, rootFilePath string) (string, error) {
	dir := filepath.Clean(start)
	for {
		if _, err := os.Stat(filepath.Join(dir, rootFilePath)); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("logger: root directory not found")
		}
		dir = parent
	}
}

func (l *Logger) Error(v ...any) {
	text := fmt.Sprintf("ERROR: %s", fmt.Sprint(v...))
	l.log.Println(text)
}

func (l *Logger) Info(v ...any) {
	text := fmt.Sprintf("INFO: %s", fmt.Sprint(v...))
	l.log.Println(text)
}

func (l *Logger) Debug(v ...any) {
	text := fmt.Sprintf("DEBUG: %s", fmt.Sprint(v...))
	l.log.Println(text)
}

func (l *Logger) Warn(v ...any) {
	text := fmt.Sprintf("WARN: %s", fmt.Sprint(v...))
	l.log.Println(text)
}

func (l *Logger) Fatal(v ...any) {
	text := fmt.Sprintf("FATAL: %s", fmt.Sprint(v...))
	l.log.Fatal(text)
}

func NewLogger(category string) (*Logger, error) {
	fileName := fmt.Sprintf("%s_%s.log", category, time.Now().Format("2006-01-02"))

	rootDir, err := findProjectRoot(".", "")
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(filepath.Join(rootDir, "log")); os.IsNotExist(err) {
		os.Mkdir(filepath.Join(rootDir, "log"), 0755)
	}

	logFilePath := filepath.Join(rootDir, "log", fileName)
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		_, err := os.Create(logFilePath)
		if err != nil {
			return nil, err
		}
	}

	f, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &Logger{
		log: log.New(f, "", log.Ldate|log.Ltime|log.Lmicroseconds),
	}, nil
}
