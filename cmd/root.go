// Package cmd is main commands packages
package cmd

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var (
	cfgFile  string
	logLevel string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "miniscrape",
	Short: "A miniscrape is an application to scrape the applications",
	Long:  `A miniscrape is an application to scrape the applications`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	var defaultLogLevel = "info"
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		defaultLogLevel = envLogLevel
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config-file", "",
		"config file (default is ./config/food.yml)")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log", "L",
		defaultLogLevel, "Set log level")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	logLevelZero, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		logLevelZero = zerolog.WarnLevel
	}

	zerolog.SetGlobalLevel(logLevelZero)

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err != nil {
			log.Err(err).
				Str("config_file", cfgFile).
				Msg("unable to load config")
		}
	} else {
		viper.AddConfigPath("./config/")

		loadConfig("default-config")
		loadConfig("local-config")
	}

	viper.AutomaticEnv() // read in environment variables that match
}

func loadConfig(name string) {
	viper.SetConfigName(name)
	if err := viper.MergeInConfig(); err != nil {
		log.Debug().
			Str("config_name", name).
			Err(err).
			Msg("unable to load config")
	}
}
