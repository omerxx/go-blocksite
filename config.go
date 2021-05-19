package main

import (
	"fmt"

	"github.com/spf13/viper"
)

func readConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func getBlackListConfig() []string {
	return viper.GetStringSlice("blacklist")
}
