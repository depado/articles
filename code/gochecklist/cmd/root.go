package cmd

import (
	"github.com/onrik/logrus/filename"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "myprogram",
	Short: "Myprogram does stuff I guess",
	Long:  "Simple program that does stuff.",
}

// Execute executes the commands
func Execute(b, v string) {
	Build = b
	Version = v
	rootCmd.AddCommand(version)
	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Fatal()
	}
}

func init() {
	cobra.OnInitialize(initialize)

	// Global flags
	rootCmd.PersistentFlags().String("conf", "", "path to the configuration file")
	rootCmd.PersistentFlags().String("log.level", "info", "one of debug, info, warn, error or fatal")
	rootCmd.PersistentFlags().String("log.format", "text", "one of text or json")
	rootCmd.PersistentFlags().Bool("log.line", false, "enable filename and line in logs")

	// Flag binding
	viper.BindPFlags(rootCmd.PersistentFlags())
}

func initialize() {
	// Environment variables
	viper.AutomaticEnv()

	if viper.GetString("conf") != "" {
		viper.SetConfigFile(viper.GetString("conf"))
	} else {
		viper.SetConfigName("conf")
		viper.AddConfigPath(".")
		viper.AddConfigPath("/config/")
	}

	// Configuration file
	if err := viper.ReadInConfig(); err != nil {
		logrus.Warn("No configuration file found")
	}

	lvl := viper.GetString("log.level")
	l, err := logrus.ParseLevel(lvl)
	if err != nil {
		logrus.WithField("level", lvl).Warn("Invalid log level, fallback to 'info'")
	} else {
		logrus.SetLevel(l)
	}
	switch viper.GetString("log.format") {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	default:
	case "text":
		logrus.SetFormatter(&logrus.TextFormatter{})
	}
	if viper.GetBool("log.line") {
		logrus.AddHook(filename.NewHook())
	}
}
