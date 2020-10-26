package main

import (
	"os"
)

func updateMain() {
	devMessage("Entering updateMain()")
	sanitizeArguments()

	updateErr := storeOnlineDatabase(config.Update.DatabaseURL, config.DatabaseFile)
	if updateErr != nil {
		stdErr.Printf("Error updating local database: %s\n", updateErr)
		os.Exit(errDatabaseUpdate)
	}

	devMessage("Leaving updateMain()")
}
