package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	. "github.com/omerxx/go-blocksite/types"
	"github.com/spf13/viper"
)

func readState() (state State) {
	stateFile := viper.GetString("app.state")
	jsonFile, err := os.Open(stateFile)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &state)
	return
}

func RemoveUrlFromList(list []string, url string) ([]string, error) {
	for index, element := range list {
		if element == url {
			return append(list[:index], list[index+1:]...), nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Couldnt find %s in %+v", url, list))
}

func FilterStateListedSites(sites []string) (filteredSites []string) {
	filteredSites = append(filteredSites, sites...)
	stateSites := readState().Blacklist
	for _, site := range sites {
		for _, stateSite := range stateSites {
			if site == stateSite.Url {
				filteredSites, _ = RemoveUrlFromList(filteredSites, stateSite.Url)
				break
			}
		}
	}
	return
}

func isListed(list []Site, target string) bool {
	for _, site := range list {
		if site.Url == target {
			return true
		}
	}
	return false
}

func AddMultiple(targets []string) {
	stateFile := viper.GetString("app.state")
	jsonFile, err := os.Open(stateFile)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	var state State
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &state)

	for _, target := range targets {
		if !isListed(state.Blacklist, target) {
			site := Site{Url: target}
			state.Blacklist = append(state.Blacklist, site)
		}
	}

	file, _ := json.MarshalIndent(state, "", " ")
	_ = ioutil.WriteFile(stateFile, file, 0644)
}

func RemoveMultiple(targets []string) {
	for _, target := range targets {
		remove(target)
	}
}

func remove(target string) {
	stateFile := viper.GetString("app.state")
	jsonFile, err := os.Open(stateFile)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	var state State
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &state)

	for i, site := range state.Blacklist {
		if site.Url == target {
			state.Blacklist = append(state.Blacklist[:i], state.Blacklist[i+1:]...)
		}
	}

	file, _ := json.MarshalIndent(state, "", " ")
	_ = ioutil.WriteFile(stateFile, file, 0644)
}

func ListSites() []Site {
	stateFile := viper.GetString("app.state")
	jsonFile, err := os.Open(stateFile)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	var state State
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &state)
	return state.Blacklist
}

func ListSitesAsStrings() (sites []string) {
	stateFile := viper.GetString("app.state")
	jsonFile, err := os.Open(stateFile)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	var state State
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &state)
	for _, site := range state.Blacklist {
		sites = append(sites, site.Url)
	}
	return sites
}
