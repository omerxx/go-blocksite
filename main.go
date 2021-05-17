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

func getBlackListSites(hosts *txeh.Hosts) []string {
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

func cleanBlocks(hosts *txeh.Hosts, sites []string) {
	for _, stateSite := range state.ListSites() {
		exists, _, _ := hosts.HostAddressLookup(stateSite.Url)
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
	sites := getBlackListSites(hosts)
	blockSites(hosts, sites)
	state.AddMultiple(sites)
	cleanBlocks(hosts, sites)
}
