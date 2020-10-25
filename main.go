package main

/*
#### ##     ## ########   #######  ########  ########  ######
 ##  ###   ### ##     ## ##     ## ##     ##    ##    ##    ##
 ##  #### #### ##     ## ##     ## ##     ##    ##    ##
 ##  ## ### ## ########  ##     ## ########     ##     ######
 ##  ##     ## ##        ##     ## ##   ##      ##          ##
 ##  ##     ## ##        ##     ## ##    ##     ##    ##    ##
#### ##     ## ##         #######  ##     ##    ##     ######
*/

import (
	"fmt"
	"log"
	"os"

	cobra "github.com/spf13/cobra"
	envordef "gitlab.com/rbrt-weiler/go-module-envordef"
)

/*
######## ##    ## ########  ########  ######
   ##     ##  ##  ##     ## ##       ##    ##
   ##      ####   ##     ## ##       ##
   ##       ##    ########  ######    ######
   ##       ##    ##        ##             ##
   ##       ##    ##        ##       ##    ##
   ##       ##    ##        ########  ######
*/

type appConfig struct {
	DatabaseFile string
	Lookup       struct {
		FetchOnlineDatabase bool
		OutputJSON          bool
	}
	Update struct {
		DatabaseURL string
	}
}

/*
 ######   #######  ##    ##  ######  ########    ###    ##    ## ########  ######
##    ## ##     ## ###   ## ##    ##    ##      ## ##   ###   ##    ##    ##    ##
##       ##     ## ####  ## ##          ##     ##   ##  ####  ##    ##    ##
##       ##     ## ## ## ##  ######     ##    ##     ## ## ## ##    ##     ######
##       ##     ## ##  ####       ##    ##    ######### ##  ####    ##          ##
##    ## ##     ## ##   ### ##    ##    ##    ##     ## ##   ###    ##    ##    ##
 ######   #######  ##    ##  ######     ##    ##     ## ##    ##    ##     ######
*/

const (
	// General constants
	toolName    string = "OUI Lookup"
	toolVersion string = "dev"
	toolID      string = toolName + "/" + toolVersion
	toolURL     string = "https://gitlab.com/rbrt-weiler/ouilookup"

	// Error codes
	errSuccess     int = 0
	errMissingArgs int = 1

	// Hardcoded defaults as fallbacks
	ouiDatabaseFile string = "oui.txt"
	ouiDatabaseURL  string = "http://standards-oui.ieee.org/oui.txt"
)

/*
##     ##    ###    ########  ####    ###    ########  ##       ########  ######
##     ##   ## ##   ##     ##  ##    ## ##   ##     ## ##       ##       ##    ##
##     ##  ##   ##  ##     ##  ##   ##   ##  ##     ## ##       ##       ##
##     ## ##     ## ########   ##  ##     ## ########  ##       ######    ######
 ##   ##  ######### ##   ##    ##  ######### ##     ## ##       ##             ##
  ## ##   ##     ## ##    ##   ##  ##     ## ##     ## ##       ##       ##    ##
   ###    ##     ## ##     ## #### ##     ## ########  ######## ########  ######
*/

var (
	config appConfig
	stdErr = log.New(os.Stderr, "", 0)
)

/*
##     ##    ###    #### ##    ##
###   ###   ## ##    ##  ###   ##
#### ####  ##   ##   ##  ####  ##
## ### ## ##     ##  ##  ## ## ##
##     ## #########  ##  ##  ####
##     ## ##     ##  ##  ##   ###
##     ## ##     ## #### ##    ##
*/

func main() {
	var rootCmd = &cobra.Command{Use: "ouilookup"}
	rootCmd.Version = toolVersion
	rootCmd.SetVersionTemplate(fmt.Sprintf("%s\n", toolID))
	rootCmd.PersistentFlags().StringVarP(&config.DatabaseFile, "dbfile", "d", envordef.StringVal("OUILOOKUP_DBFILE", "oui.txt"), "Local database file to use")

	var cmdLookup = &cobra.Command{
		Use:   "lookup [stuff]",
		Short: "Look up MAC vendor",
		Long:  `Some longer description for lookup here.`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("looking up MAC vendor")
		},
	}
	cmdLookup.Flags().BoolVarP(&config.Lookup.FetchOnlineDatabase, "dbfetch", "f", envordef.BoolVal("OUILOOKUP_FETCHDB", true), "Fetch online database if offline not accessible")
	cmdLookup.Flags().BoolVarP(&config.Lookup.OutputJSON, "json", "j", envordef.BoolVal("OUILOOKUP_JSON", false), "Output in JSON format")

	var cmdUpdate = &cobra.Command{
		Use:   "update [stuff]",
		Short: "Update OUI database",
		Long:  `Some longer description for update here.`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("updating database")
		},
	}
	cmdUpdate.Flags().StringVarP(&config.Update.DatabaseURL, "dburl", "u", envordef.StringVal("OUILOOKUP_DBURL", ouiDatabaseURL), "URL to fetch the database from")

	rootCmd.AddCommand(cmdLookup, cmdUpdate)
	rootCmd.Execute()

	fmt.Printf("%v\n", config)

	os.Exit(errSuccess)
}
