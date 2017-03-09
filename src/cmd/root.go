package cmd

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"github.com/cloudposse/openvpn-api/src/server"
	"github.com/cloudposse/openvpn-api/src/config"
)

var cfgFile string

var flags = []flag{
	{"l", "string", "listen", ":8085", "Listen              ( environment variable LISTEN could be used instead )"},
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "openvpn-api",
	Short: "Use OpenVPN API to create users and fetch predefined openvpn configs",
	Long: `
Use OpenVPN API to create users and fetch predefined openvpn configs.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := log.WithFields(log.Fields{"class": "RootCmd", "method": "RunE"})

		cfg := config.Config{
			Listen: viper.GetString("listen"),
		}

		logger.Infof("Config: Listen - %v", cfg.Listen)

		err := cfg.Validate()

		if err == nil {
			server.Run(cfg)
		}

		return err
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Config file
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"Config file         (default is $HOME/.openvpn-api.yaml)")

	for _, f := range flags {
		createCmdFlags(RootCmd, f)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".openvpn-api") // name of config file (without extension)
	viper.AddConfigPath("$HOME")                   // adding home directory as first search path
	viper.AutomaticEnv()                           // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
