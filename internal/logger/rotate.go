package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// DailyRotatingWriter writes to a new file each calendar day.
// Files are named <dir>/<prefix>-YYYY-MM-DD.log.
// Thread-safe; safe to use as a logrus output writer.
type DailyRotatingWriter struct {
	mu     sync.Mutex
	dir    string
	prefix string
	maxAge int // keep this many days of log files; 0 = keep forever
	date   string
	file   *os.File
}

// Option configures a DailyRotatingWriter.
type Option func(*DailyRotatingWriter)

// WithMaxAge sets how many days of log files to retain.
// Files older than maxAge days are deleted on each rotation.
// A value of 0 disables pruning.
func WithMaxAge(days int) Option {
	return func(w *DailyRotatingWriter) {
		w.maxAge = days
	}
}

// NewDailyRotatingWriter creates and opens a new DailyRotatingWriter.
// It creates dir if it does not exist.
func NewDailyRotatingWriter(dir, prefix string, opts ...Option) (*DailyRotatingWriter, error) {
	w := &DailyRotatingWriter{
		dir:    dir,
		prefix: prefix,
		maxAge: 30,
	}
	for _, o := range opts {
		o(w)
	}
	if err := w.openForToday(); err != nil {
		return nil, err
	}
	return w, nil
}

// Write implements io.Writer. It rotates the file if the calendar day has changed.
func (w *DailyRotatingWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if err := w.rotateIfNeeded(); err != nil {
		return 0, err
	}
	return w.file.Write(p)
}

// Sync flushes buffered data to the OS.
func (w *DailyRotatingWriter) Sync() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.file != nil {
		return w.file.Sync()
	}
	return nil
}

// Close flushes and closes the current log file.
func (w *DailyRotatingWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.file != nil {
		_ = w.file.Sync()
		return w.file.Close()
	}
	return nil
}

func (w *DailyRotatingWriter) rotateIfNeeded() error {
	today := currentDate()
	if w.date == today {
		return nil
	}
	if w.file != nil {
		_ = w.file.Sync()
		_ = w.file.Close()
		w.file = nil
	}
	if err := w.openForToday(); err != nil {
		return err
	}
	if w.maxAge > 0 {
		w.pruneOldFiles()
	}
	return nil
}

func (w *DailyRotatingWriter) openForToday() error {
	if err := os.MkdirAll(w.dir, 0o755); err != nil {
		return fmt.Errorf("logger: create log dir %q: %w", w.dir, err)
	}
	today := currentDate()
	path := filepath.Join(w.dir, fmt.Sprintf("%s-%s.log", w.prefix, today))
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("logger: open log file %q: %w", path, err)
	}
	w.file = f
	w.date = today
	return nil
}

// pruneOldFiles removes log files older than w.maxAge days.
// Called without the lock held — only invoked from rotateIfNeeded which holds it.
func (w *DailyRotatingWriter) pruneOldFiles() {
	cutoff := time.Now().AddDate(0, 0, -w.maxAge)
	pattern := filepath.Join(w.dir, w.prefix+"-*.log")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}
	for _, path := range matches {
		name := filepath.Base(path)
		// strip prefix and ".log" to get the date portion
		dateStr := strings.TrimPrefix(name, w.prefix+"-")
		dateStr = strings.TrimSuffix(dateStr, ".log")
		t, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}
		if t.Before(cutoff) {
			_ = os.Remove(path)
		}
	}
}

func currentDate() string {
	return time.Now().Format("2006-01-02")
}
