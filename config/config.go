package config

import (
  "os"
	"fmt"
  "errors"
  "golang.org/x/sys/unix"
	"github.com/spf13/viper"
	flag "github.com/spf13/pflag"
)

type Options struct {
  GlobalUnblock bool
  Block         string
  Unblock       string
}

type Config struct {
  blockRoute string
  hostsfile   string
  state       string
  blockList  []string
}

func setConfigDefaults() (config Config) {
  config = Config{
    blockRoute: "localhost",
    hostsfile: "/etc/hosts",
    state: "state.json",
  }
  return config
}

func ReadConfig() (config Config){
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
    config = setConfigDefaults()
		// panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
  config.hostsfile = viper.GetString("app.hostsfile")
  config.blockRoute = viper.GetString("app.blockRoute")
  config.state = viper.GetString("app.state")
  config.blockList = viper.GetStringSlice("blacklist")
  return config
}

func addWsiteUrls(sites []string) (updatedList []string) {
  updatedList = append(updatedList, sites...)
  for _, site := range sites {
    wsite := fmt.Sprintf("www.%s", site)
    updatedList = append(updatedList, wsite)
  }
	return
}

func getBlackList(config Config) (list []string){
  list = config.blockList
  list = addWsiteUrls(list)
	return
}

func ParseOptions() Options {
	var globalUnblock *bool = flag.BoolP("unblock-all", "U", false, "unblocks all urls")
	var blockUrl *string = flag.StringP("block", "b", "", "blocks a given url")
	var unblockUrl *string = flag.StringP("unblock", "u", "", "unblocks a given url")
  flag.Parse()
  options := Options{
    GlobalUnblock: *globalUnblock,
    Block: *blockUrl,
    Unblock: *unblockUrl,
  }
  return options
}

func addUrlToList(list []string, url string) []string {
  return append(list, url)
}

func removeUrlFromList(list []string, url string) ([]string, error) {
  for index, element := range list {
    if element == url {
      return append(list[:index], list[index+1:]...), nil
    }
  }
  return nil, errors.New(fmt.Sprintf("Couldnt find %s in %+v", url, list))
}

func wsite(site string) (wsite string) {
  return fmt.Sprintf("www.%s", site)
}

func updateBlackListWithCli(configuredBlackList []string, options Options) (updatedList []string) {
  updatedList = append(updatedList, configuredBlackList...)
  if options.Block != "" {
    updatedList = addUrlToList(updatedList, options.Block)
    updatedList = addUrlToList(updatedList, wsite(options.Block))
  }
  if options.Unblock != "" {
    var err error
    updatedList, err = removeUrlFromList(updatedList, options.Unblock)
    if err !=nil {
      fmt.Print(err)
    }
    updatedList, err = removeUrlFromList(updatedList, wsite(options.Unblock))
    if err !=nil {
      fmt.Print(err)
    }
  }
  return
}

func writeable(path string) bool {
    return unix.Access(path, unix.W_OK) == nil
    // return true
}

func canWriteToHosts(config Config) bool {
  hostsfile := config.hostsfile
  return writeable(hostsfile)
}

func RunPreflightChecks(config Config) {
  if !canWriteToHosts(config) {
    fmt.Println("Can't write to hostsfile, aborting.")
    os.Exit(0)
  }
}

func HandleOptions(blacklistConfiguredSites []string, config Config) []string {
  options := ParseOptions()
  if !options.GlobalUnblock {
    blacklistConfiguredSites = getBlackList(config)
    blacklistConfiguredSites = updateBlackListWithCli(blacklistConfiguredSites, options)
  }
  return blacklistConfiguredSites
}
