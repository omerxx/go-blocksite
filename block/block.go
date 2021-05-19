package block

import (
  "fmt"
	"github.com/omerxx/go-blocksite/state"
	"github.com/spf13/viper"
	"github.com/txn2/txeh"
)

func BlockSites(hosts *txeh.Hosts, sites []string) {
	for _, site := range sites {
		blockSite(hosts, site)
	}
}

func blockSite(hosts *txeh.Hosts, site string) {
	target := viper.GetString("app.blockTarget")

	hosts.AddHost(target, site)
  wsite := fmt.Sprintf("www.%s", site)
	hosts.AddHost(target, wsite)
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

func CleanBlocks(hosts *txeh.Hosts, blacklistConfiguredSites []string) (cleanedTargets []string) {
	stateSites := state.ListSites()
	for _, stateSite := range stateSites {
		exists := isInList(stateSite.Url, blacklistConfiguredSites)
		if !exists {
			hosts.RemoveHost(stateSite.Url)
			hosts.Save()
			// state.Remove(stateSite.Url)
			cleanedTargets = append(cleanedTargets, stateSite.Url)
		}
	}
	return
}
