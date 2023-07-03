package logger

import (
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	logruspapertrail "github.com/polds/logrus-papertrail-hook"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
	"golang.org/x/sys/unix"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Log is the global logger
var Log *logrus.Entry
var BaseLogger *logrus.Logger

// LoggingConfig represents the configuration for the logging mechanism of the application.
type LoggingConfig struct {
	Debug        bool     // if true, log level will be set to DEBUG
	Logfile      string   // empty to disable logfile
	MaxSize      int      // rotate after x MB filesize
	MaxAge       int      // rotate evey x days
	Backup       int      // number of backups to keep
	PTLogLevel   string   // <DEBUG|INFO|WARNING|ERROR> or empty to disable PT logging
	PTAppName    string   // if empty name of executable will be used
	PTSourceHost string   // Which host should be shown in PaperTrail
	PTHost       string   // PaperTrail host
	PTPort       int      // PaperTrail port
	FieldsOrder  []string // order of fields in log output
}

// Isatty Check if a TTY is attached
func Isatty() bool {
	_, err := unix.IoctlGetWinsize(int(os.Stdout.Fd()), unix.TIOCGWINSZ)
	if err != nil {
		return false
	}
	return true
}

// Trace returns a string with useful trace information used in error handling
func Trace() string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	funcName := strings.SplitN(frame.Function, "/", 2)
	return fmt.Sprintf("%s %s:%d", funcName[len(funcName)-1], filepath.Base(frame.File), frame.Line)
}

// InitLogger Initialize a new logger and set it as the global logger
func InitLogger(cfg LoggingConfig) {
	BaseLogger = newLogger(cfg)
	Log = BaseLogger.WithFields(logrus.Fields{})
}

// newLogger Initialize a new logger
func newLogger(cfg LoggingConfig) *logrus.Logger {
	newLogger := logrus.New()
	newLogger.SetOutput(io.Discard)
	if cfg.Debug {
		newLogger.SetLevel(logrus.DebugLevel)
	} else {
		newLogger.SetLevel(logrus.InfoLevel)
	}

	// Set defaults
	if cfg.MaxSize == 0 {
		cfg.MaxSize = 10
	}

	if cfg.MaxAge == 0 {
		cfg.MaxAge = 30
	}

	if cfg.Backup == 0 {
		cfg.Backup = 6
	}

	// ########################################################################
	// Output to console
	// ########################################################################
	if Isatty() {
		newLogger.Hooks.Add(lfshook.NewHook(
			os.Stdout,
			&nested.Formatter{
				NoColors:         false,
				HideKeys:         true,
				NoFieldsColors:   false,
				NoFieldsSpace:    false,
				ShowFullLevel:    false,
				NoUppercaseLevel: false,
				TrimMessages:     true,
				TimestampFormat:  "15:04:05",
				FieldsOrder:      cfg.FieldsOrder,
			},
		))
	}

	// ########################################################################
	// Output to logfile
	// ########################################################################
	if cfg.Logfile != "" {
		rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
			Filename:   cfg.Logfile,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.Backup,
			MaxAge:     cfg.MaxAge,
			Level:      logrus.DebugLevel,
			Formatter: &nested.Formatter{
				HideKeys:         true,
				NoColors:         true,
				NoFieldsColors:   false,
				NoFieldsSpace:    false,
				ShowFullLevel:    false,
				NoUppercaseLevel: false,
				TrimMessages:     true,
				TimestampFormat:  "2006-01-02 15:04:05",
				FieldsOrder:      cfg.FieldsOrder,
			},
		})
		if err != nil {
			logrus.Fatalf("Failed to initialize logfile: %v", err)
		}
		newLogger.AddHook(rotateFileHook)
	}

	// ########################################################################
	// Output to Papertrail
	// ########################################################################
	if cfg.PTLogLevel != "" {
		newLogger.SetFormatter(&logrus.TextFormatter{
			DisableTimestamp: true,
			ForceColors:      false,
		})

		var appName string
		if cfg.PTAppName == "" {
			filename, _ := os.Executable()
			appName = filepath.Base(filename)
		} else {
			appName = cfg.PTAppName
		}

		var hook *logruspapertrail.Hook
		hook, err := logruspapertrail.NewPapertrailHook(&logruspapertrail.Hook{
			Host:     "logs2.papertrailapp.com",
			Port:     39413,
			Hostname: "bigdata01",
			Appname:  appName,
		})

		switch strings.ToUpper(cfg.PTLogLevel) {
		case "DEBUG":
			hook.SetLevels([]logrus.Level{logrus.DebugLevel, logrus.InfoLevel, logrus.WarnLevel, logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel})
		case "INFO":
			hook.SetLevels([]logrus.Level{logrus.InfoLevel, logrus.WarnLevel, logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel})
		case "WARNING":
			hook.SetLevels([]logrus.Level{logrus.WarnLevel, logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel})
		case "ERROR":
			hook.SetLevels([]logrus.Level{logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel})
		default:
			hook.SetLevels([]logrus.Level{})
		}

		if err == nil {
			newLogger.AddHook(hook)
		}
	}

	return newLogger
}
