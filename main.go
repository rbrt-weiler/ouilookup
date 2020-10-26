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
	Update       struct {
		DatabaseURL        string
		HTTPTimeoutSeconds uint
	}
	Export struct {
		OutputFormat string
	}
	Lookup struct {
		FetchOnlineDatabase bool
		OutputJSON          bool
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
	errSuccess        int = 0
	errMissingArgs    int = 1
	errDatabaseUpdate int = 10
	errDatabaseLoad   int = 15
	errDatabaseParse  int = 16
	errExportFormat   int = 20

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
		stdErr.Printf("DEV %v: %s\n", time.Now(), message)
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
	devMessage("Entering main()")

	var rootCmd = &cobra.Command{Use: "ouilookup"}
	rootCmd.Version = toolVersion
	rootCmd.SetVersionTemplate(fmt.Sprintf("%s\n", toolID))
	rootCmd.PersistentFlags().StringVarP(&config.DatabaseFile, "dbfile", "d", envordef.StringVal("OUILOOKUP_DBFILE", ouiDatabaseFile), "Local database file to use")

	var cmdUpdate = &cobra.Command{
		Use:   "update",
		Short: "Update local OUI database",
		Long:  `Some longer description for update here.`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			updateMain()
		},
	}
	cmdUpdate.Flags().StringVarP(&config.Update.DatabaseURL, "dburl", "u", envordef.StringVal("OUILOOKUP_DBURL", ouiDatabaseURL), "URL to fetch the database from")
	cmdUpdate.Flags().UintVarP(&config.Update.HTTPTimeoutSeconds, "httptimeout", "t", envordef.UintVal("OUILOOKUP_HTTPTIMEOUT", 60), "HTTP timeout in seconds")

	var cmdExport = &cobra.Command{
		Use:   "export",
		Short: "Export OUI database",
		Long:  `Some longer description for export here.`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			exportMain()
		},
	}
	cmdExport.Flags().StringVar(&config.Export.OutputFormat, "format", envordef.StringVal("OUILOOKUP_EXPORTFORMAT", "csv"), "Output format for export")

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

	rootCmd.AddCommand(cmdUpdate, cmdExport, cmdLookup)
	rootCmd.Execute()

	devMessage("Leaving main()")
	os.Exit(errSuccess)
}
