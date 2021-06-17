package config

import (
  "os"
  "fmt"
  "errors"
  "golang.org/x/sys/unix"
  "github.com/spf13/viper"
  flag "github.com/spf13/pflag"
	"github.com/omerxx/go-blocksite/state"
	. "github.com/omerxx/go-blocksite/types"
)

type Options struct {
  GlobalUnblock bool
  Blocks        []string
  Unblocks      []string
  List          bool
}

type Config struct {
  blockRoute string
  hostsfile  string
  state      string
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

// func getConifg(configEntry string) string {
//   if viper.GetString(configEntry) == "" {
//     return defaultConfigs.con
//   }
// }

func ReadConfig() (config Config){
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
    config = setConfigDefaults()
	}
  config.hostsfile = viper.GetString("app.hostsfile")
  config.blockRoute = viper.GetString("app.blockRoute")
  config.state = viper.GetString("app.state")
  config.blockList = append(config.blockList, readBlockList()...)

  return config
}

func readBlockList() []string {
  configSitesEnabled := viper.GetBool("blocklist.enabled")
  if configSitesEnabled {
    return viper.GetStringSlice("blocklist.sites")
  }
  return []string{}
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
  var blockUrls *[]string = flag.StringSliceP("block", "b", []string{}, "blocks a given url(s)")
	var unblockUrls *[]string = flag.StringSliceP("unblock", "u", []string{}, "unblocks given url(s)")
	var list *bool = flag.BoolP("list", "l", false, "lists current blocked sites")
  flag.Parse()
  options := Options{
    GlobalUnblock: *globalUnblock,
    Blocks: *blockUrls,
    Unblocks: *unblockUrls,
    List: *list,
  }
  return options
}

func addUrlsToList(list []string, urls []string) []string {
  return append(list, urls...)
}

func removeUrlFromList(list []string, url string) ([]string, error) {
  for index, element := range list {
    if element == url {
      return append(list[:index], list[index+1:]...), nil
    }
  }
  return nil, errors.New(fmt.Sprintf("Couldnt find %s in %+v", url, list))
}

// TODO: There's a bug here where when removing a non-existing URL it removes the entire state
func removeUrlsFromList(list []string, urls []string) (updatedList []string, err error) {
  updatedList = append(updatedList, list...)
  for _, url := range urls {
    updatedList, err = removeUrlFromList(list, url)
  }
  return updatedList, err
}

func wsites(sites []string) (wsites []string) {
  for _, site := range sites {
    wsites = append(wsites, fmt.Sprintf("www.%s", site))
  }
  return wsites
}

func updateBlackListWithCli(configuredBlackList []string, options Options) (updatedList []string) {
  updatedList = append(updatedList, configuredBlackList...)
  updatedList = append(updatedList, state.ListSitesAsStrings()...)
  // if options.Block != "" {
  if len(options.Blocks) != 0 {
    updatedList = addUrlsToList(updatedList, options.Blocks)
    updatedList = addUrlsToList(updatedList, wsites(options.Blocks))
  }
  if len(options.Unblocks) != 0 {
    var err error
    updatedList, err = removeUrlsFromList(updatedList, options.Unblocks)
    if err !=nil {
      fmt.Print(err)
    }
    updatedList, err = removeUrlsFromList(updatedList, wsites(options.Unblocks))
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

func prettyPrintListedSites(sites []Site) {
  for _, site := range sites {
    fmt.Println(site.Url)
  }
}

func HandleOptions(blacklistConfiguredSites []string, config Config) []string {
  options := ParseOptions()
  if !options.GlobalUnblock {
    blacklistConfiguredSites = getBlackList(config)
    blacklistConfiguredSites = updateBlackListWithCli(blacklistConfiguredSites, options)
  }
  if options.List {
    prettyPrintListedSites(state.ListSites())
    os.Exit(0)
  }
  return blacklistConfiguredSites
}
