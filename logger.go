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

func LogDbg(metadata ContextfulReqMetadata, message string) {
	logEnvConfig, logEnvConfigExists := os.LookupEnv(GO_BLAST_LOG_LEVEL)
	if !logEnvConfigExists {
		logAny(GO_BLAST_LOG_DEBUG_LEVEL, metadata.ReferenceId, metadata.TracingId, message)
	} else {
		if prohibitLogTypeOnLogLevels([]string{
			GO_BLAST_LOG_INFO_LEVEL,
			GO_BLAST_LOG_WARN_LEVEL,
			GO_BLAST_LOG_ERROR_LEVEL,
		}, logEnvConfig) {
			logAny(GO_BLAST_LOG_DEBUG_LEVEL, metadata.ReferenceId, metadata.TracingId, message)
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

func LogInf(metadata ContextfulReqMetadata, message string) {
	logEnvConfig, logEnvConfigExists := os.LookupEnv(GO_BLAST_LOG_LEVEL)
	if !logEnvConfigExists {
		logAny(GO_BLAST_LOG_INFO_LEVEL, metadata.ReferenceId, metadata.TracingId, message)
	} else {
		if prohibitLogTypeOnLogLevels([]string{
			GO_BLAST_LOG_WARN_LEVEL,
			GO_BLAST_LOG_ERROR_LEVEL,
		}, logEnvConfig) {
			logAny(GO_BLAST_LOG_INFO_LEVEL, metadata.ReferenceId, metadata.TracingId, message)
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

func LogWrn(metadata ContextfulReqMetadata, message string) {
	logEnvConfig, logEnvConfigExists := os.LookupEnv(GO_BLAST_LOG_LEVEL)
	if !logEnvConfigExists {
		logAny(GO_BLAST_LOG_WARN_LEVEL, metadata.ReferenceId, metadata.TracingId, message)
	} else {
		if prohibitLogTypeOnLogLevels([]string{
			GO_BLAST_LOG_ERROR_LEVEL,
		}, logEnvConfig) {
			logAny(GO_BLAST_LOG_WARN_LEVEL, metadata.ReferenceId, metadata.TracingId, message)
		}
	}
}

func LogError(referenceId, tracingId, message string) {
	logAny(GO_BLAST_LOG_ERROR_LEVEL, referenceId, tracingId, message)
}

func LogErr(metadata ContextfulReqMetadata, message string) {
	logAny(GO_BLAST_LOG_ERROR_LEVEL, metadata.ReferenceId, metadata.TracingId, message)
}
