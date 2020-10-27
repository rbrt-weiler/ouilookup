package main

import (
	"fmt"
	"os"
)

func vendorMain(args []string) {
	devMessage("Entering vendorMain()")

	db, dbErr := loadDatabase(config.DatabaseFile)
	if dbErr != nil {
		stdErr.Printf("Error loading database: %s\n", dbErr)
		os.Exit(errDatabaseLoad)
	}
	vendorDB, vendorDBErr := ouiToVendorDatabase(db)
	if vendorDBErr != nil {
		stdErr.Printf("Error converting database: %s\n", vendorDBErr)
		os.Exit(errDatabaseConvert)
	}

	for _, vendor := range args {
		if ouis, vendorExists := vendorDB[vendor]; vendorExists {
			for _, oui := range ouis {
				mac, macErr := normalizeMAC(oui)
				if macErr != nil {
					stdErr.Printf("Warning: MAC could not be normalized: %s\n", macErr)
					continue
				}
				fmt.Printf("%s = %s\n", vendor, mac[:8])
			}
		}
	}

	devMessage("Leaving vendorMain()")
}
