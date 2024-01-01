package cmd

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"

	// Enable logrus logging in baaah
	_ "github.com/acorn-io/baaah/pkg/logrus"
	sloglogrus "github.com/samber/slog-logrus"

	"github.com/acorn-io/cmd/pkg/logserver"
	"github.com/google/go-containerregistry/pkg/logs"
	"github.com/sirupsen/logrus"
	"k8s.io/klog/v2"
)

type DebugLogging struct {
	Debug                 bool `usage:"Enable debug logging"`
	DebugLevel            int  `usage:"Debug log level (valid 0-9) (default 7)"`
	EnableDynamicLogLevel bool `usage:"Enable loglevel server to enable changing the log level at runtime"`
}

func setIfUnset(envKey, value string) {
	if v := os.Getenv(envKey); v == "" {
		os.Setenv(envKey, value)
	}
}

func (d DebugLogging) InitLogging() error {
	slog.SetDefault(slog.New(sloglogrus.Option{
		Level: slog.LevelInfo,
	}.NewLogrusHandler()))

	if level := os.Getenv("ACORN_LOG_LEVEL"); level != "" && !d.Debug && d.DebugLevel == 0 {
		switch level {
		case "trace":
			d.DebugLevel = 7
		case "debug":
			d.DebugLevel = 6
		}
	}

	if d.Debug || d.DebugLevel > 0 {
		slog.SetDefault(slog.New(sloglogrus.Option{
			Level: slog.LevelDebug,
		}.NewLogrusHandler()))

		logging := flag.NewFlagSet("", flag.PanicOnError)
		klog.InitFlags(logging)

		level := d.DebugLevel
		if level == 0 {
			level = 6
		}
		if level > 7 {
			setIfUnset("ACORN_LOG_LEVEL", "trace")
			logrus.SetLevel(logrus.TraceLevel)
			logs.Debug = log.New(os.Stderr, "ggcr: ", log.LstdFlags)
		} else {
			setIfUnset("ACORN_LOG_LEVEL", "debug")
			logrus.SetLevel(logrus.DebugLevel)
		}
		if err := logging.Parse([]string{
			fmt.Sprintf("-v=%d", level),
		}); err != nil {
			return err
		}
		logrus.Debug("Debug logging enabled")
	}

	if d.EnableDynamicLogLevel {
		logrus.Debug("Log server enabled")
		go logserver.StartServerWithDefaults()
	}

	return nil
}
