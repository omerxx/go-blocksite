package main

import (
	"github.com/txn2/txeh"
	"github.com/omerxx/go-blocksite/block"
	"github.com/omerxx/go-blocksite/state"
	"github.com/omerxx/go-blocksite/config"
)

// TODO
func backupHostsFile() {}

func readConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

// TODO
func canWriteToHosts() bool {
	return false
}

func getBlackListConfig() []string {
	return viper.GetStringSlice("blacklist")
}

func blockSites(hosts *txeh.Hosts, sites []string) {
	for _, site := range sites {
		blockSite(hosts, site)
	}
}

func blockSite(hosts *txeh.Hosts, site string) {
	target := viper.GetString("app.blockTarget")

	hosts.AddHost(target, site)
	hosts.Save()
}

func isInList(object string, list []string) bool {
	for _, element := range list {
		if element == object {
			return true
		}
	}
	return false
}

func cleanBlocks(hosts *txeh.Hosts) {
	blacklistConfiguredSites := getBlackListConfig()
	stateSites := state.ListSites()
	for _, stateSite := range stateSites {
		exists := isInList(stateSite.Url, blacklistConfiguredSites)
		if !exists {
			hosts.RemoveHost(stateSite.Url)
			hosts.Save()
			state.Remove(stateSite.Url)
		}
	}
}

func main() {
  configuration := config.ReadConfig()
  config.RunPreflightChecks(configuration)
	hosts, err := txeh.NewHostsDefault()
	if err != nil {
		panic(err)
	}
  blacklistConfiguredSites := []string{}
  blacklistConfiguredSites = config.HandleOptions(blacklistConfiguredSites, configuration)
	blacklistConfiguredSites = block.BlockSites(hosts, blacklistConfiguredSites)

	state.AddMultiple(blacklistConfiguredSites)
	cleanedTargets := block.CleanBlocks(hosts, blacklistConfiguredSites)
	state.RemoveMultiple(cleanedTargets)
}
