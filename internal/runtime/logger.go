package runtime

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"container-hooks-toolkit/internal/logger"

	log "github.com/sirupsen/logrus"
)

// Logger adds a way to manage output to a log file to a logrus.Logger
type Logger struct {
	logger.Interface
	previousLogger logger.Interface
	logFiles       []*os.File
}

// NewLogger creates an empty logger
func NewLogger() *Logger {
	return &Logger{
		Interface: log.New(),
	}
}

// Update constructs a Logger with a preddefined formatter
func (l *Logger) Update(filename string, logLevel string, argv []string) {

	configFromArgs := parseArgs(argv)

	level, logLevelError := configFromArgs.getLevel(logLevel)

	// If the logLevel is not specified, we use infoLevel, and log warning it after update
	// which means the progress won't terminated if the logLevel is not specified
	defer func() {
		if logLevelError != nil {
			l.Warning(logLevelError)
		}
	}()

	var logFiles []*os.File
	var argLogFileError error

	// We don't create log files if the version argument is supplied
	if !configFromArgs.version {
		configLogFile, err := logger.CreateLogFile(filename)
		if err != nil {
			argLogFileError = errors.Join(argLogFileError, err)
		}
		if configLogFile != nil {
			logFiles = append(logFiles, configLogFile)
		}

		argLogFile, err := logger.CreateLogFile(configFromArgs.file)
		if argLogFile != nil {
			logFiles = append(logFiles, argLogFile)
		}
		if err != nil {
			argLogFileError = errors.Join(argLogFileError, err)
		}
	}
	defer func() {
		if argLogFileError != nil {
			l.Warningf("Failed to open log file: %v", argLogFileError)
		}
	}()

	newLogger := log.New()

	newLogger.SetLevel(level)
	if level == log.DebugLevel {
		log.SetReportCaller(true)
		// Shorten function and file names reported by the logger, by
		// trimming common "github.com/opencontainers/runc" prefix.
		// This is only done for text formatter.
		_, file, _, _ := runtime.Caller(0)
		prefix := filepath.Dir(file) + "/"
		log.SetFormatter(&log.TextFormatter{
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				function := strings.TrimPrefix(f.Function, prefix) + "()"
				fileLine := strings.TrimPrefix(f.File, prefix) + ":" + strconv.Itoa(f.Line)
				return function, fileLine
			},
		})
	}

	if configFromArgs.format == "json" {
		newLogger.SetFormatter(new(log.JSONFormatter))
	}

	if len(logFiles) == 0 {
		newLogger.SetOutput(io.Discard)
	} else if len(logFiles) == 1 {
		newLogger.SetOutput(logFiles[0])
	} else if len(logFiles) > 1 {
		var writers []io.Writer
		for _, f := range logFiles {
			writers = append(writers, f)
		}
		newLogger.SetOutput(io.MultiWriter(writers...))
	}

	*l = Logger{
		Interface:      newLogger,
		previousLogger: l.Interface,
		logFiles:       logFiles,
	}
}

// Reset closes the log file (if any) and resets the logger output to what it
// was before UpdateLogger was called.
func (l *Logger) Reset() error {
	defer func() {
		previous := l.previousLogger
		if previous == nil {
			previous = log.New()
		}
		l.Interface = previous
		l.previousLogger = nil
		l.logFiles = nil
	}()

	var errs []error
	for _, f := range l.logFiles {
		err := f.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}

	var err error
	if len(errs) > 0 {
		return fmt.Errorf("%v", errs)
	}

	return err
}

type loggerConfig struct {
	file    string
	format  string
	debug   bool
	version bool
}

func (c loggerConfig) getLevel(logLevel string) (log.Level, error) {
	if c.debug {
		return log.DebugLevel, nil
	}

	if logLevel, err := log.ParseLevel(logLevel); err == nil {
		return logLevel, nil
	}

	return log.InfoLevel, fmt.Errorf("invalid log-level '%v'", logLevel)
}

func parseArgs(args []string) loggerConfig {
	c := loggerConfig{}

	expected := map[string]*string{
		"log-format": &c.format,
		"log":        &c.file,
	}

	found := make(map[string]bool)

	for i := 0; i < len(args); i++ {
		if len(found) == 4 {
			break
		}

		param := args[i]

		parts := strings.SplitN(param, "=", 2)
		trimmed := strings.TrimLeft(parts[0], "-")
		// If this is not a flag we continue
		if parts[0] == trimmed {
			continue
		}

		// Check the version flag
		if trimmed == "version" {
			c.version = true
			found["version"] = true
			// For the version flag we don't process any other flags
			continue
		}

		// Check the debug flag
		if trimmed == "debug" {
			c.debug = true
			found["debug"] = true
			continue
		}

		destination, exists := expected[trimmed]
		if !exists {
			continue
		}

		var value string
		if len(parts) == 2 {
			value = parts[2]
		} else if i+1 < len(args) {
			value = args[i+1]
			i++
		} else {
			continue
		}

		*destination = value
		found[trimmed] = true
	}

	return c
}
