/*
	Copyright © 2022 Peter Stanko <peter.stanko0@gmail.com>
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var (
	cfgFile    string
	logEnabled bool
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
	rootCmd.PersistentFlags().BoolVarP(&logEnabled, "log", "L",
		false, "Enable logging")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if !logEnabled {
		log.SetOutput(ioutil.Discard)
	}
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
		viper.ReadInConfig()

	} else {
		viper.AddConfigPath("./config/")

		loadConfig("default-config")
		loadConfig("local-config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.MergeInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func loadConfig(name string) {
	viper.SetConfigName(name)
	log.Printf("Load config file: %s\n", viper.ConfigFileUsed())
	viper.MergeInConfig()
}
