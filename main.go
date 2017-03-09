package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/cloudposse/openvpn-api/src/cmd"
	"github.com/spf13/viper"
	"os"
)

func main() {
	LoggerInit()
	cmd.Execute()
}

// LoggerInit - Initialize logger configuration used for cli
func LoggerInit() {
	viper.SetDefault("log_level", "info")
	viper.BindEnv("log_level", "LOG_LEVEL")

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(os.Stderr)

	// Only log the warning severity or above.
	loglevel := viper.GetString("log_level")
	switch loglevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
		break

	case "info":
	default:
		log.SetLevel(log.InfoLevel)
		break
	}
}