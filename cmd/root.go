package cmd

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"os"

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
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config-file", "",
		"config file (default is ./config/food.yml)")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log", "L",
		"info", "Set log level")
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
			log.Error().
				Str("config_file", cfgFile).
				Err(err).
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
		log.Error().
			Str("config_name", name).
			Err(err).
			Msg("unable to load config")
	}
}
