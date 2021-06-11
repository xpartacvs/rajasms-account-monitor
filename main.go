package main

import (
	"rajasms-account-monitor/packages/logger"
	"rajasms-account-monitor/packages/worker"
)

func main() {
	logger.Log().Info().Msg("Application is starting...")

	if err := worker.Start(); err != nil {
		logger.Log().Err(err).Msg("Application has been terminated due to unexpected error.")
	}
}
