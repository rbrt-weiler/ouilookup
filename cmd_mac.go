package main

import (
	"fmt"
	"os"
	"strings"
)

func macMain(args []string) {
	var vendorName string

	devMessage("Entering macMain()")
	sanitizeArguments()

	db, dbErr := loadDatabase(config.DatabaseFile)
	if dbErr != nil {
		stdErr.Printf("Error loading database: %s\n", dbErr)
		os.Exit(errDatabaseLoad)
	}

	for _, mac := range args {
		if !isValidMAC(mac) {
			stdErr.Printf("Warning: MAC %s is invalid.\n", mac)
			continue
		}
		mac = strings.ToLower(mac)
		oui, ouiErr := extractOUI(mac)
		if ouiErr != nil {
			stdErr.Printf("Warning: Unable to map OUI for MAC %s: %s\n", mac, ouiErr)
			continue
		}
		if vendorName = db.OUIDatabase[oui].VendorName; vendorName == "" {
			vendorName = "(unregistered)"
		}
		normalized, normalizedErr := normalizeMAC(mac)
		if normalizedErr != nil {
			stdErr.Printf("Warning: MAC could not be normalized: %s\n", normalizedErr)
			continue
		}
		fmt.Printf("%s = %s\n", normalized, vendorName)
	}

	devMessage("Leaving macMain()")
}
