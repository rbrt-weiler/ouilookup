package main

import (
	"fmt"
	"os"
)

func exportMain() {
	var output string

	devMessage("Entering exportMain()")
	sanitizeArguments()

	db, dbErr := loadDatabase(config.DatabaseFile)
	if dbErr != nil {
		stdErr.Printf("Error reading local OUI database: %s\n", dbErr)
		os.Exit(errDatabaseLoad)
	}

	switch config.Export.OutputFormat {
	case "text":
		output = db.ToText()
	case "csv":
		output = db.ToCSV()
	case "json":
		output = db.ToJSON()
	default:
		stdErr.Printf("Error: Unsupported export format.")
		os.Exit(errExportFormat)
	}

	fmt.Println(output)
}
