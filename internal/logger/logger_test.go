package logger

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestDailyRotatingWriter_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	w, err := NewDailyRotatingWriter(dir, "test")
	if err != nil {
		t.Fatalf("NewDailyRotatingWriter: %v", err)
	}
	defer w.Close()

	if _, err := w.Write([]byte("hello\n")); err != nil {
		t.Fatalf("Write: %v", err)
	}

	files, _ := filepath.Glob(filepath.Join(dir, "test-*.log"))
	if len(files) != 1 {
		t.Fatalf("expected 1 log file, got %d", len(files))
	}
	content, _ := os.ReadFile(files[0])
	if string(content) != "hello\n" {
		t.Fatalf("unexpected content: %q", content)
	}
}

func TestWriterHook_WritesJSONToFile(t *testing.T) {
	dir := t.TempDir()
	w, err := NewDailyRotatingWriter(dir, "app")
	if err != nil {
		t.Fatalf("NewDailyRotatingWriter: %v", err)
	}
	defer w.Close()

	log := logrus.New()
	log.SetOutput(io.Discard)
	log.SetLevel(logrus.DebugLevel)
	log.AddHook(&writerHook{
		writer:    w,
		formatter: &logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"},
		levels:    logrus.AllLevels,
	})

	log.WithField("key", "value").Info("hook test")

	_ = w.Sync()

	files, _ := filepath.Glob(filepath.Join(dir, "app-*.log"))
	if len(files) != 1 {
		t.Fatalf("expected 1 log file, got %d", len(files))
	}
	content, _ := os.ReadFile(files[0])
	t.Logf("file content: %s", content)

	var entry map[string]interface{}
	if err := json.Unmarshal(content, &entry); err != nil {
		t.Fatalf("log file is not valid JSON: %v\ncontent: %s", err, content)
	}
	if entry["msg"] != "hook test" {
		t.Fatalf("expected msg=%q, got %v", "hook test", entry["msg"])
	}
}

func TestCreateLoggerInstant_WritesToFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("LOG_DIR", dir)
	t.Setenv("HEX_ARCH_ENV", "")

	log := CreateLoggerInstant()
	log.SetOutput(io.Discard)
	log.Info("integration test message")

	today := currentDate()
	expectedFile := filepath.Join(dir, "app-"+today+".log")
	content, err := os.ReadFile(expectedFile)
	if err != nil {
		t.Fatalf("log file not found at %s: %v", expectedFile, err)
	}
	if !strings.Contains(string(content), "integration test message") {
		t.Fatalf("message not found in log file.\ncontent: %s", content)
	}
}

// TestCreateLoggerInstant_ProductionWritesToFile verifies that even in production
// mode (HEX_ARCH_ENV=release, console level=Error), Info messages still reach the file.
func TestCreateLoggerInstant_ProductionWritesToFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("LOG_DIR", dir)
	t.Setenv("HEX_ARCH_ENV", "release")

	log := CreateLoggerInstant()
	log.Info("should appear in file despite production mode")

	today := currentDate()
	expectedFile := filepath.Join(dir, "app-"+today+".log")
	content, err := os.ReadFile(expectedFile)
	if err != nil {
		t.Fatalf("log file not found at %s: %v", expectedFile, err)
	}
	if !strings.Contains(string(content), "should appear in file despite production mode") {
		t.Fatalf("Info message missing from file in production mode.\ncontent: %s", content)
	}
}
