package goblast

import (
	"log"
	"os"
	"slices"
)

const (
	GO_BLAST_LOG_LEVEL = "GOBLAST-LOG-LEVEL"

	GO_BLAST_LOG_ERROR_LEVEL = "ERROR"
	GO_BLAST_LOG_WARN_LEVEL  = "WARN"
	GO_BLAST_LOG_INFO_LEVEL  = "INFO"
	GO_BLAST_LOG_DEBUG_LEVEL = "DEBUG"
)

// Base Logger method
func logAny(logType, referenceId, tracingId, message string) {
	logPrefix := "[" + referenceId + "," + tracingId + "][" + logType + "] - "
	log.Println(logPrefix + message)
}

// To filter out prohibited log levels from performing logging
func prohibitLogTypeOnLogLevels(prohibitedLevels []string, logLevel string) bool {
	return !slices.Contains(prohibitedLevels, logLevel)
}

func LogDebug(referenceId, tracingId, message string) {
	logEnvConfig, logEnvConfigExists := os.LookupEnv(GO_BLAST_LOG_LEVEL)
	if !logEnvConfigExists {
		logAny(GO_BLAST_LOG_DEBUG_LEVEL, referenceId, tracingId, message)
	} else {
		if prohibitLogTypeOnLogLevels([]string{
			GO_BLAST_LOG_INFO_LEVEL,
			GO_BLAST_LOG_WARN_LEVEL,
			GO_BLAST_LOG_ERROR_LEVEL,
		}, logEnvConfig) {
			logAny(GO_BLAST_LOG_DEBUG_LEVEL, referenceId, tracingId, message)
		}
	}
}

func LogInfo(referenceId, tracingId, message string) {
	logEnvConfig, logEnvConfigExists := os.LookupEnv(GO_BLAST_LOG_LEVEL)
	if !logEnvConfigExists {
		logAny(GO_BLAST_LOG_INFO_LEVEL, referenceId, tracingId, message)
	} else {
		if prohibitLogTypeOnLogLevels([]string{
			GO_BLAST_LOG_WARN_LEVEL,
			GO_BLAST_LOG_ERROR_LEVEL,
		}, logEnvConfig) {
			logAny(GO_BLAST_LOG_INFO_LEVEL, referenceId, tracingId, message)
		}
	}
}

func LogWarn(referenceId, tracingId, message string) {
	logEnvConfig, logEnvConfigExists := os.LookupEnv(GO_BLAST_LOG_LEVEL)
	if !logEnvConfigExists {
		logAny(GO_BLAST_LOG_WARN_LEVEL, referenceId, tracingId, message)
	} else {
		if prohibitLogTypeOnLogLevels([]string{
			GO_BLAST_LOG_ERROR_LEVEL,
		}, logEnvConfig) {
			logAny(GO_BLAST_LOG_WARN_LEVEL, referenceId, tracingId, message)
		}
	}
}

func LogError(referenceId, tracingId, message string) {
	logAny(GO_BLAST_LOG_ERROR_LEVEL, referenceId, tracingId, message)
}
