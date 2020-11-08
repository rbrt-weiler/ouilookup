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
	Server struct {
		HTTPPort uint
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
	toolVersion string = "0.3.0"
	toolID      string = toolName + "/" + toolVersion
	toolURL     string = "https://gitlab.com/rbrt-weiler/ouilookup"

	// Error codes
	errSuccess         int = 0
	errMissingArgs     int = 1
	errDatabaseUpdate  int = 10
	errDatabaseLoad    int = 15
	errDatabaseParse   int = 16
	errDatabaseConvert int = 17
	errExportFormat    int = 20

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
	if config.Server.HTTPPort < 1024 || config.Server.HTTPPort > 65535 {
		config.Server.HTTPPort = 8000
	}
	devMessage("Leaving sanitizeArguments()")
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
		Long:  `Use update to fetch a copy of an online OUI database and save it locally.`,
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
		Long: `Use export to export the locally stored OUI database in various formats.
Valid output formats are "text", "csv" and "json". Output is written to stdout.`,
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			exportMain()
		},
	}
	cmdExport.Flags().StringVarP(&config.Export.OutputFormat, "format", "f", envordef.StringVal("OUILOOKUP_EXPORTFORMAT", "csv"), "Output format for export")

	var cmdMAC = &cobra.Command{
		Use:   "mac [mac...]",
		Short: "Look up MAC vendor",
		Long:  `Use mac to retrieve the vendor name of any number of given MAC addresses.`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			macMain(args)
		},
	}

	var cmdVendor = &cobra.Command{
		Use:   "vendor [name...]",
		Short: "Look up OUIs for a specific vendor",
		Long:  `Use vendor to retrieve the OUIs assigned to a specific vendor.`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			vendorMain(args)
		},
	}

	var cmdServer = &cobra.Command{
		Use:   "server",
		Short: "Start a HTTP server to lookup MACs and vendors",
		Long: `Use server to start a HTTP server. The server provides a RESTful API that
can be used to lookup MACs and vendors.`,
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			serverMain()
		},
	}
	cmdServer.Flags().UintVar(&config.Server.HTTPPort, "port", envordef.UintVal("OUILOOKUP_HTTP_PORT", 8000), "HTTP port to listen on")

	rootCmd.AddCommand(cmdUpdate, cmdExport, cmdMAC, cmdVendor, cmdServer)
	rootCmd.Execute()

	devMessage("Leaving main()")
	os.Exit(errSuccess)
}
