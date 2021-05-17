package main

import (
	"fmt"

	"github.com/omerxx/go-blocksite/state"

	"github.com/spf13/viper"
	"github.com/txn2/txeh"
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
		fmt.Printf("statesite %s in configured is %t\n", stateSite.Url, exists)
		// exists, _, _ := hosts.HostAddressLookup(stateSite.Url)
		if !exists {
			hosts.RemoveHost(stateSite.Url)
			hosts.Save()
			state.Remove(stateSite.Url)
		}
	}
}

func main() {
	readConfig()
	hosts, err := txeh.NewHostsDefault()
	if err != nil {
		panic(err)
	}
	blacklistConfiguredSites := getBlackListConfig()
	blockSites(hosts, blacklistConfiguredSites)
	state.AddMultiple(blacklistConfiguredSites)
	cleanBlocks(hosts)
}
