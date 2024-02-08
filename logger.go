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
func logAny(logType, referenceId, message string) {
	logPrefix := "[" + referenceId + "][" + logType + "] - "
	log.Println(logPrefix + message)
}

// To filter out prohibited log levels from performing logging
func prohibitLogTypeOnLogLevels(prohibitedLevels []string, logLevel string) bool {
	return !slices.Contains(prohibitedLevels, logLevel)
}

func LogDebug(referenceId, message string) {
	logEnvConfig, logEnvConfigExists := os.LookupEnv(GO_BLAST_LOG_LEVEL)
	if !logEnvConfigExists {
		logAny(GO_BLAST_LOG_DEBUG_LEVEL, referenceId, message)
	} else {
		if prohibitLogTypeOnLogLevels([]string{
			GO_BLAST_LOG_INFO_LEVEL,
			GO_BLAST_LOG_WARN_LEVEL,
			GO_BLAST_LOG_ERROR_LEVEL,
		}, logEnvConfig) {
			logAny(GO_BLAST_LOG_DEBUG_LEVEL, referenceId, message)
		}
	}
}

func LogDbg(metadata ContextfulReqMetadata, message string) {
	logEnvConfig, logEnvConfigExists := os.LookupEnv(GO_BLAST_LOG_LEVEL)
	if !logEnvConfigExists {
		logAny(GO_BLAST_LOG_DEBUG_LEVEL, metadata.ReferenceId, message)
	} else {
		if prohibitLogTypeOnLogLevels([]string{
			GO_BLAST_LOG_INFO_LEVEL,
			GO_BLAST_LOG_WARN_LEVEL,
			GO_BLAST_LOG_ERROR_LEVEL,
		}, logEnvConfig) {
			logAny(GO_BLAST_LOG_DEBUG_LEVEL, metadata.ReferenceId, message)
		}
	}
}

func LogInfo(referenceId, message string) {
	logEnvConfig, logEnvConfigExists := os.LookupEnv(GO_BLAST_LOG_LEVEL)
	if !logEnvConfigExists {
		logAny(GO_BLAST_LOG_INFO_LEVEL, referenceId, message)
	} else {
		if prohibitLogTypeOnLogLevels([]string{
			GO_BLAST_LOG_WARN_LEVEL,
			GO_BLAST_LOG_ERROR_LEVEL,
		}, logEnvConfig) {
			logAny(GO_BLAST_LOG_INFO_LEVEL, referenceId, message)
		}
	}
}

func LogInf(metadata ContextfulReqMetadata, message string) {
	logEnvConfig, logEnvConfigExists := os.LookupEnv(GO_BLAST_LOG_LEVEL)
	if !logEnvConfigExists {
		logAny(GO_BLAST_LOG_INFO_LEVEL, metadata.ReferenceId, message)
	} else {
		if prohibitLogTypeOnLogLevels([]string{
			GO_BLAST_LOG_WARN_LEVEL,
			GO_BLAST_LOG_ERROR_LEVEL,
		}, logEnvConfig) {
			logAny(GO_BLAST_LOG_INFO_LEVEL, metadata.ReferenceId, message)
		}
	}
}

func LogWarn(referenceId, message string) {
	logEnvConfig, logEnvConfigExists := os.LookupEnv(GO_BLAST_LOG_LEVEL)
	if !logEnvConfigExists {
		logAny(GO_BLAST_LOG_WARN_LEVEL, referenceId, message)
	} else {
		if prohibitLogTypeOnLogLevels([]string{
			GO_BLAST_LOG_ERROR_LEVEL,
		}, logEnvConfig) {
			logAny(GO_BLAST_LOG_WARN_LEVEL, referenceId, message)
		}
	}
}

func LogWrn(metadata ContextfulReqMetadata, message string) {
	logEnvConfig, logEnvConfigExists := os.LookupEnv(GO_BLAST_LOG_LEVEL)
	if !logEnvConfigExists {
		logAny(GO_BLAST_LOG_WARN_LEVEL, metadata.ReferenceId, message)
	} else {
		if prohibitLogTypeOnLogLevels([]string{
			GO_BLAST_LOG_ERROR_LEVEL,
		}, logEnvConfig) {
			logAny(GO_BLAST_LOG_WARN_LEVEL, metadata.ReferenceId, message)
		}
	}
}

func LogError(referenceId, message string) {
	logAny(GO_BLAST_LOG_ERROR_LEVEL, referenceId, message)
}

func LogErr(metadata ContextfulReqMetadata, message string) {
	logAny(GO_BLAST_LOG_ERROR_LEVEL, metadata.ReferenceId, message)
}
