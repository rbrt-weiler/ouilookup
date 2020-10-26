package main

import (
	"fmt"
	"os"
	"time"

	resty "github.com/go-resty/resty/v2"
)

func updateMain() {
	devMessage("Entering updateMain()")
	sanitizeArguments()

	onlineDatabase, onlineDatabaseErr := fetchOnlineDatabase()
	if onlineDatabaseErr != nil {
		stdErr.Printf("Error fetching online OUI database: %s\n", onlineDatabaseErr)
		os.Exit(errDatabaseFetch)
	}

	compressedData, compressedDataErr := compressData(onlineDatabase)
	if compressedDataErr != nil {
		stdErr.Printf("Error compressing OUI database: %s\n", compressedDataErr)
		os.Exit(errDatabaseCompress)
	}

	storeErr := storeData(config.DatabaseFile, compressedData.String())
	if storeErr != nil {
		stdErr.Printf("Error storing local OUI database: %\n", storeErr)
		os.Exit(errDatabaseStore)
	}
}

func fetchOnlineDatabase() (string, error) {
	devMessage("Entering fetchOnlineDatabase()")

	client := resty.New()
	client.SetTimeout(time.Duration(config.Update.HTTPTimeoutSeconds) * time.Second)
	resp, respErr := client.R().Get(config.Update.DatabaseURL)
	if respErr != nil {
		return "", respErr
	}
	devMessage(fmt.Sprintf("Status Code   : %v", resp.StatusCode()))
	devMessage(fmt.Sprintf("Response Size : %v", resp.Size()))

	return resp.String(), nil
}
