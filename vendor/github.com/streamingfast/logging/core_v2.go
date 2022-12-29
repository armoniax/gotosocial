package logging

import (
	"fmt"
	"os"

	"github.com/blendle/zapdriver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/ssh/terminal"
)

// This v2 version of `core.go` is a work in progress without any backwar compatibility
// version. It might not made it to an official version of the library so you can depend
// at your own risk.

type loggerOptions struct {
	autoStartSwitcherServer *bool
	encoderVerbosity        *int
	level                   *zap.AtomicLevel
	loggerName              *string
	reportAllErrors         *bool
	serviceName             *string
	zapOptions              []zap.Option
}

type LoggerOption interface {
	apply(o *loggerOptions)
}

type loggerFuncOption func(o *loggerOptions)

func (f loggerFuncOption) apply(o *loggerOptions) {
	f(o)
}

func WithAutoStartSwitcherServer() LoggerOption {
	return loggerFuncOption(func(o *loggerOptions) {
		o.autoStartSwitcherServer = ptrBool(true)
	})
}

func WithAtomicLevel(level zap.AtomicLevel) LoggerOption {
	return loggerFuncOption(func(o *loggerOptions) {
		o.level = ptrLevel(level)
	})
}

func WithLoggerName(name string) LoggerOption {
	return loggerFuncOption(func(o *loggerOptions) {
		o.loggerName = ptrString(name)
	})
}

func WithReportAllErrors() LoggerOption {
	return loggerFuncOption(func(o *loggerOptions) {
		o.reportAllErrors = ptrBool(true)
	})
}

func WithServiceName(name string) LoggerOption {
	return loggerFuncOption(func(o *loggerOptions) {
		o.serviceName = ptrString(name)
	})
}

func WithZapOption(zapOption zap.Option) LoggerOption {
	return loggerFuncOption(func(o *loggerOptions) {
		o.zapOptions = append(o.zapOptions, zapOption)
	})
}

// LibraryLogger creates a new no-op logger (via `zap.NewNop`) and automatically registered it
// withing the logging registry with a tracer that can be be used for conditionally tracing
// code.
func LibraryLogger(shortName string, packageID string, logger **zap.Logger) Tracer {
	return libraryLogger(globalRegistry, shortName, packageID, logger)
}

func libraryLogger(registry *registry, shortName string, packageID string, logger **zap.Logger) Tracer {
	return register2(registry, shortName, packageID, logger)
}

// ApplicationLogger should be used to get a logger for a top-level binary application.
//
// By default,
func ApplicationLogger(shortName string, packageID string, logger **zap.Logger, opts ...LoggerOption) Tracer {
	return applicationLogger(globalRegistry, os.Getenv, shortName, packageID, logger, opts...)
}

func applicationLogger(
	registry *registry,
	envGet func(string) string,
	shortName string,
	packageID string,
	logger **zap.Logger,
	opts ...LoggerOption,
) Tracer {
	loggerOptions := loggerOptions{}
	for _, opt := range opts {
		opt.apply(&loggerOptions)
	}

	if loggerOptions.reportAllErrors == nil {
		WithReportAllErrors().apply(&loggerOptions)
	}

	if loggerOptions.serviceName == nil {
		WithServiceName(shortName).apply(&loggerOptions)
	}

	if loggerOptions.autoStartSwitcherServer == nil && isProductionEnvironment() {
		opts = append(opts, WithAutoStartSwitcherServer())
	}

	tracer := register2(registry, shortName, packageID, logger)

	loggerFactory := func(level zapcore.Level) *zap.Logger {
		// If the level was specified up-front, let's not use the one received
		if loggerOptions.level != nil {
			return newLogger(&loggerOptions)
		}

		clonedOptions := loggerOptions
		clonedOptions.level = ptrLevel(zap.NewAtomicLevelAt(level))

		return newLogger(&clonedOptions)
	}

	// We must keep the pointer because it could be moved in the override below
	initialLogger := *logger

	logLevelSpec := newLogLevelSpec(envGet)
	registry.overrideFromSpec(logLevelSpec, loggerFactory)

	appEntry := registry.entriesByPackageID[packageID]
	if *appEntry.logPtr != nil && *appEntry.logPtr == initialLogger {
		// No environment override the default logger, let's force INFO to be used in this case
		registry.setLoggerForEntry(appEntry, zapcore.InfoLevel, false, loggerFactory)
	}

	return tracer
}

// NewLogger creates a new logger with sane defaults based on a varity of rules described
// below and automatically registered withing the logging registry.
func NewLogger(opts ...LoggerOption) *zap.Logger {
	logger, err := MaybeNewLogger(opts...)
	if err != nil {
		panic(fmt.Errorf("unable to create logger (in production? %t): %w", isProductionEnvironment(), err))
	}

	return logger
}

func MaybeNewLogger(opts ...LoggerOption) (*zap.Logger, error) {
	options := loggerOptions{}
	for _, opt := range opts {
		opt.apply(&options)
	}

	logger, err := maybeNewLogger(&options)
	if err != nil {
		return nil, err
	}

	if options.loggerName != nil {
		logger = logger.Named(*options.loggerName)
	}

	return logger, nil
}

func newLogger(opts *loggerOptions) *zap.Logger {
	logger, err := maybeNewLogger(opts)
	if err != nil {
		panic(fmt.Errorf("unable to create logger (in production? %t): %w", isProductionEnvironment(), err))
	}

	return logger
}

func maybeNewLogger(opts *loggerOptions) (*zap.Logger, error) {
	zapOptions := opts.zapOptions

	if isProductionEnvironment() {
		reportAllErrors := opts.reportAllErrors != nil
		serviceName := opts.serviceName

		if reportAllErrors && opts.serviceName != nil {
			zapOptions = append(zapOptions, zapdriver.WrapCore(zapdriver.ReportAllErrors(true), zapdriver.ServiceName(*serviceName)))
		} else if reportAllErrors {
			zapOptions = append(zapOptions, zapdriver.WrapCore(zapdriver.ReportAllErrors(true)))
		} else if opts.serviceName != nil {
			zapOptions = append(zapOptions, zapdriver.WrapCore(zapdriver.ServiceName(*serviceName)))
		}

		return zapdriver.NewProductionConfig().Build(zapOptions...)
	}

	// Development logger
	isTTY := terminal.IsTerminal(int(os.Stderr.Fd()))
	logStdoutWriter := zapcore.Lock(os.Stderr)
	verbosity := 1
	if opts.encoderVerbosity != nil {
		verbosity = *opts.encoderVerbosity
	}

	return zap.New(zapcore.NewCore(NewEncoder(verbosity, isTTY), logStdoutWriter, opts.level)), nil
}

// func newDefaultLoggerOptions() (o *loggerOptions) {
// 	return &loggerOptions{
// 		encoderVerbosity: inferEncoderVerbosity(),
// 		level:            inferLevel(),
// 	}
// }

// func inferLevel() zap.AtomicLevel {
// 	if os.Getenv("DEBUG") != "" || os.Getenv("TRACE") != "" {
// 		return zap.NewAtomicLevelAt(zapcore.DebugLevel)
// 	}

// 	return zap.NewAtomicLevelAt(zapcore.InfoLevel)
// }

// func inferEncoderVerbosity() int {
// 	if os.Getenv("DEBUG") != "" || os.Getenv("TRACE") != "" {
// 		return 3
// 	}

// 	return 1
// }

func isProductionEnvironment() bool {
	_, err := os.Stat("/.dockerenv")

	return !os.IsNotExist(err)
}

type Tracer interface {
	Enabled() bool
}

type boolTracer struct {
	value *bool
}

func (t boolTracer) Enabled() bool {
	if t.value == nil {
		return false
	}

	return *t.value
}

func ptrBool(value bool) *bool                        { return &value }
func ptrInt(value int) *int                           { return &value }
func ptrString(value string) *string                  { return &value }
func ptrLevel(value zap.AtomicLevel) *zap.AtomicLevel { return &value }
