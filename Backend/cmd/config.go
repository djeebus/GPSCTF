package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"log"
	"os"
)

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigType("yaml")
	if cfgFile != "" {
		fmt.Println("Trying to use", cfgFile)
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Fatalf("Failed to get homedir: %v", err)
		}

		workDir, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get workdir: %v", err)
		}

		// Search config in home directory with name ".gpsctf" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(workDir)
		viper.SetConfigName(".gpsctf")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	fmt.Println("Using config file:", viper.ConfigFileUsed())
}
