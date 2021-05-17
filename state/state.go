package state

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/viper"
)

type Site struct {
	Url string `json:"url"`
}

type State struct {
	Blacklist []Site `json:"blacklist"`
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
	defer jsonFile.Close()
}

func Remove(target string) {
	stateFile := viper.GetString("app.state")
	jsonFile, err := os.Open(stateFile)
	if err != nil {
		fmt.Println(err)
	}

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
	defer jsonFile.Close()
}

func ListSites() []Site {
	stateFile := viper.GetString("app.state")
	jsonFile, err := os.Open(stateFile)
	if err != nil {
		fmt.Println(err)
	}
	var state State
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &state)
	defer jsonFile.Close()
	return state.Blacklist
}
