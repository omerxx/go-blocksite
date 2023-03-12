package main

import (
	"github.com/omerxx/go-blocksite/block"
	"github.com/omerxx/go-blocksite/config"
	"github.com/omerxx/go-blocksite/state"
	"github.com/txn2/txeh"
)

// // TODO
// func backupHostsFile() {}

// // TODO
// func createStateIfdoesntExist() {}

// // TODO
// func showBlockList() {}

// // TODO
// func unblockAllAndSave() {}

// // TODO
// func blockFromSavedFile() {}

// // TODO
// func setSchedule() {}

// // TODO
// func generateSocialBlcokList() {}

// // TODO
// func generateMediaBlockList() {}

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
