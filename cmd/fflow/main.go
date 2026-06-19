package main

import (
	"log"
	"os"

	"github.com/baibeicha/fflow/internal/fflow/cli"
	"github.com/baibeicha/fflow/pkg/logger"
)

func main() {
	appLogger, err := logger.New(
		logger.WithEnv("dev"),
		logger.WithLevel("err"),
	)
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer func(appLogger *logger.Logger) {
		err := appLogger.Close()
		if err != nil {
			log.Fatalf("failed to close logger: %v", err)
		}
	}(appLogger)

	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
