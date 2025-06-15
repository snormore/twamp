package twamp_test

import (
	"flag"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/lmittmann/tint"
)

var (
	log *slog.Logger
)

func TestMain(m *testing.M) {
	flag.Parse()
	verbose := false
	if vFlag := flag.Lookup("test.v"); vFlag != nil && vFlag.Value.String() == "true" {
		verbose = true
	}
	log = newTestLogger(verbose)

	os.Exit(m.Run())
}

func newTestLogger(verbose bool) *slog.Logger {
	logWriter := os.Stdout
	logLevel := slog.LevelDebug
	if !verbose {
		logLevel = slog.LevelInfo
	}
	logger := slog.New(tint.NewHandler(logWriter, &tint.Options{
		Level:      logLevel,
		TimeFormat: time.Kitchen,
		AddSource:  true,
	}))
	return logger
}
