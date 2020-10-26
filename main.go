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
	"time"

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
		DatabaseURL        string
		HTTPTimeoutSeconds uint
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
	// If this file exists, devMode is set to true
	devModeFile string = ".devMode"

	// General constants
	toolName    string = "OUI Lookup"
	toolVersion string = "dev"
	toolID      string = toolName + "/" + toolVersion
	toolURL     string = "https://gitlab.com/rbrt-weiler/ouilookup"

	// Error codes
	errSuccess          int = 0
	errMissingArgs      int = 1
	errDatabaseFetch    int = 10
	errDatabaseCompress int = 11
	errDatabaseStore    int = 12

	// Hardcoded defaults as fallbacks
	ouiDatabaseFile string = "oui.txt.gz"
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
	config  appConfig
	devMode bool = false
	stdErr       = log.New(os.Stderr, "", 0)
)

/*
######## ##     ## ##    ##  ######   ######
##       ##     ## ###   ## ##    ## ##    ##
##       ##     ## ####  ## ##       ##
######   ##     ## ## ## ## ##        ######
##       ##     ## ##  #### ##             ##
##       ##     ## ##   ### ##    ## ##    ##
##        #######  ##    ##  ######   ######
*/

func setDevMode() {
	info, err := os.Stat(devModeFile)
	if err == nil {
		if !info.IsDir() {
			devMode = true
		}
	}
}

func devMessage(message string) {
	if devMode {
		stdErr.Printf("DEV: %s\n", message)
	}
}

func sanitizeArguments() {
	devMessage("Entering sanitizeArguments()")
	if config.Update.HTTPTimeoutSeconds < 5 {
		config.Update.HTTPTimeoutSeconds = 5
	} else if config.Update.HTTPTimeoutSeconds > 300 {
		config.Update.HTTPTimeoutSeconds = 300
	}
}

/*
#### ##    ## #### ########
 ##  ###   ##  ##     ##
 ##  ####  ##  ##     ##
 ##  ## ## ##  ##     ##
 ##  ##  ####  ##     ##
 ##  ##   ###  ##     ##
#### ##    ## ####    ##
*/

func init() {
	setDevMode()
}

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
	devMessage(fmt.Sprintf("Started at %v", time.Now()))

	var rootCmd = &cobra.Command{Use: "ouilookup"}
	rootCmd.Version = toolVersion
	rootCmd.SetVersionTemplate(fmt.Sprintf("%s\n", toolID))
	rootCmd.PersistentFlags().StringVarP(&config.DatabaseFile, "dbfile", "d", envordef.StringVal("OUILOOKUP_DBFILE", ouiDatabaseFile), "Local database file to use")

	var cmdLookup = &cobra.Command{
		Use:   "lookup [stuff]",
		Short: "Look up MAC vendor",
		Long:  `Some longer description for lookup here.`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			lookupMain()
		},
	}
	cmdLookup.Flags().BoolVarP(&config.Lookup.FetchOnlineDatabase, "dbfetch", "f", envordef.BoolVal("OUILOOKUP_FETCHDB", true), "Fetch online database if offline not accessible")
	cmdLookup.Flags().BoolVarP(&config.Lookup.OutputJSON, "json", "j", envordef.BoolVal("OUILOOKUP_JSON", false), "Output in JSON format")

	var cmdUpdate = &cobra.Command{
		Use:   "update",
		Short: "Update OUI database",
		Long:  `Some longer description for update here.`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			updateMain()
		},
	}
	cmdUpdate.Flags().StringVarP(&config.Update.DatabaseURL, "dburl", "u", envordef.StringVal("OUILOOKUP_DBURL", ouiDatabaseURL), "URL to fetch the database from")
	cmdUpdate.Flags().UintVarP(&config.Update.HTTPTimeoutSeconds, "httptimeout", "t", envordef.UintVal("OUILOOKUP_HTTPTIMEOUT", 60), "HTTP timeout in seconds")

	rootCmd.AddCommand(cmdLookup, cmdUpdate)
	rootCmd.Execute()

	devMessage(fmt.Sprintf("Exiting at %v", time.Now()))
	os.Exit(errSuccess)
}
