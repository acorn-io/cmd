package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"

	// Enable logrus logging in baaah
	_ "github.com/acorn-io/baaah/pkg/logrus"

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

func (d DebugLogging) InitLogging() error {
	if d.Debug || d.DebugLevel > 0 {
		logging := flag.NewFlagSet("", flag.PanicOnError)
		klog.InitFlags(logging)

		level := d.DebugLevel
		if level == 0 {
			level = 6
		}
		if level > 7 {
			logrus.SetLevel(logrus.TraceLevel)
			logs.Debug = log.New(os.Stderr, "ggcr: ", log.LstdFlags)
		} else {
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
