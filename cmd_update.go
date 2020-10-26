package main

import (
	"bytes"
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

	storeErr := storeData(config.DatabaseFile, compressedData)
	if storeErr != nil {
		stdErr.Printf("Error storing local OUI database: %\n", storeErr)
		os.Exit(errDatabaseStore)
	}

	devMessage("Leaving updateMain()")
}

func fetchOnlineDatabase() (bytes.Buffer, error) {
	var buf bytes.Buffer

	devMessage("Entering fetchOnlineDatabase()")

	client := resty.New()
	client.SetTimeout(time.Duration(config.Update.HTTPTimeoutSeconds) * time.Second)
	resp, respErr := client.R().Get(config.Update.DatabaseURL)
	if respErr != nil {
		return buf, respErr
	}
	devMessage(fmt.Sprintf("Status Code   : %v", resp.StatusCode()))
	devMessage(fmt.Sprintf("Response Size : %v", resp.Size()))

	devMessage("Leaving fetchOnlineDatabase()")
	return *bytes.NewBuffer(resp.Body()), nil
}
