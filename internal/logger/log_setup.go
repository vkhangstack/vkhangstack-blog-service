package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func logDir() string {
	if d := os.Getenv("LOG_DIR"); d != "" {
		return d
	}
	return "logs"
}

func logMaxAge() int {
	if s := os.Getenv("LOG_MAX_AGE"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n >= 0 {
			return n
		}
	}
	return 30
}

func GetLoggingEnv() string {
	if os.Getenv("HEX_ARCH_ENV") == "release" {
		return "structured"
	}
	return "stdout"
}

func SetupLogger() {
	Log = CreateLoggerInstant()
}

func CreateLoggerInstant() *logrus.Logger {
	log := logrus.New()
	log.SetReportCaller(true)
	// Always DebugLevel so every entry reaches the hooks.
	// Per-output level filtering is done inside each hook's Levels().
	log.SetLevel(logrus.DebugLevel)

	jsonFmt := &logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"}

	if GetLoggingEnv() == "structured" {
		// Production: discard default output; add a hook that sends Error+ to stdout.
		log.SetOutput(io.Discard)
		log.SetFormatter(jsonFmt)
		log.AddHook(&writerHook{
			writer:    os.Stdout,
			formatter: jsonFmt,
			levels:    errorAndAbove,
		})
	} else {
		log.SetOutput(os.Stdout)
		log.SetFormatter(&myFormatter{logrus.TextFormatter{
			FullTimestamp:          true,
			TimestampFormat:        "2006-01-02 15:04:05",
			ForceColors:            true,
			DisableLevelTruncation: true,
		}})
	}

	dir, _ := filepath.Abs(logDir())
	fileWriter, err := NewDailyRotatingWriter(dir, "app", WithMaxAge(logMaxAge()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "logger: file logging disabled (%v)\n", err)
	} else {
		fmt.Fprintf(os.Stdout, "logger: writing to %s/app-%s.log (max age %d days)\n",
			dir, currentDate(), logMaxAge())
		log.AddHook(&writerHook{
			writer:    fileWriter,
			formatter: jsonFmt,
			levels:    logrus.AllLevels,
		})
	}

	return log
}

// errorAndAbove are the levels written to stdout in production mode.
var errorAndAbove = []logrus.Level{
	logrus.PanicLevel,
	logrus.FatalLevel,
	logrus.ErrorLevel,
}

// writerHook writes formatted log entries to any io.Writer at the specified levels.
type writerHook struct {
	writer    io.Writer
	formatter logrus.Formatter
	levels    []logrus.Level
}

func (h *writerHook) Levels() []logrus.Level { return h.levels }

func (h *writerHook) Fire(entry *logrus.Entry) error {
	b, err := h.formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = h.writer.Write(b)
	return err
}

type myFormatter struct {
	logrus.TextFormatter
}

func (f *myFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var levelColor int
	strList := strings.Split(entry.Caller.File, "/")
	fileName := strList[len(strList)-1]

	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = 31 // gray
	case logrus.WarnLevel:
		levelColor = 33 // yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = 31 // red
	default:
		levelColor = 36 // blue
	}
	return []byte(fmt.Sprintf("[%s] - %s - [line:%d] - \x1b[%dm%s\x1b[0m - %s. Data: %v\n",
		entry.Time.Format(f.TimestampFormat), fileName, entry.Caller.Line, levelColor,
		strings.ToUpper(entry.Level.String()), entry.Message, entry.Data)), nil
}
