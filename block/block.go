package block

import (
	"fmt"
	"net/http"

	"github.com/omerxx/go-blocksite/state"
	"github.com/spf13/viper"
	"github.com/txn2/txeh"
)

func isTargetAlive(url string) bool {
	resp, err := http.Get(fmt.Sprintf("https://%s", url))
	if err != nil {
		// fmt.Printf("Can't probe: %s\n", url)
		// TODO switch back --->
		// TODO
		return true
		// return false
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 500 {
		return true
	} else {
		fmt.Printf("Status code for %s is %d, Ignoring\n", url, resp.StatusCode)
		return false
	}
}

func BlockSites(hosts *txeh.Hosts, sites []string) (blacklistConfiguredSites []string) {
	for _, s := range sites {
		blacklistConfiguredSites = append(blacklistConfiguredSites, s)
	}
	sites = state.FilterStateListedSites(sites)

	for _, site := range sites {
		if isTargetAlive(site) {
			blockSite(hosts, site)
		} else {
			blacklistConfiguredSites, _ = state.RemoveUrlFromList(blacklistConfiguredSites, site)
		}
	}
	return blacklistConfiguredSites
}

func blockSite(hosts *txeh.Hosts, site string) {
	target := viper.GetString("app.blockRoute")
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
