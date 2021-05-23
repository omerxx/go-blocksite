package config

import (
	"fmt"

	"github.com/spf13/viper"
	flag "github.com/spf13/pflag"
)

type Options struct {
  globalUnblock bool
  block         string
  unblock       string
}

func readConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func add_wsite_urls(sites []string) (updatedList []string) {
  updatedList = append(updatedList, sites...)
  for _, site := range sites {
    wsite := fmt.Sprintf("www.%s", site)
    updatedList = append(updatedList, wsite)
  }
	return
}

func getBlackList() (list []string){
  list = viper.GetStringSlice("blacklist")
  list = add_wsite_urls(list)
	return
}

func parseOptions() Options {
	var globalUnblock *bool = flag.BoolP("unblock-all", "U", false, "unblocks all urls")
	var blockUrl *string = flag.StringP("block", "b", "", "blocks a given url")
	var unblockUrl *string = flag.StringP("unblock", "u", "", "unblocks a given url")
  flag.Parse()
  options := Options{
    globalUnblock: *globalUnblock,
    block: *blockUrl,
    unblock: *unblockUrl,
  }
  return options
}

func updateBlackListWithCli(configuredBlackList []string, options Options) (updatedList []string) {
  updatedList = append(updatedList, configuredBlackList...)
  if options.block != "" {
    updatedList = addUrlToList(updatedList, options.block)
    updatedList = addUrlToList(updatedList, wsite(options.block))
  }
  if options.unblock != "" {
    var err error
    updatedList, err = removeUrlFromList(updatedList, options.unblock)
    if err !=nil {
      fmt.Print(err)
    }
    updatedList, err = removeUrlFromList(updatedList, wsite(options.unblock))
    if err !=nil {
      fmt.Print(err)
    }
  }
  return
}

