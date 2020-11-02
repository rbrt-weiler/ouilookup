package main

import (
	"fmt"
	"net/http"
	"os"

	mux "github.com/gorilla/mux"
)

var (
	persistentOUIDatabase    ouiDatabase
	persistentVendorDatabase map[string][]string
)

func serverMain() {
	devMessage("Entering serverMain()")
	sanitizeArguments()

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
	persistentOUIDatabase = db
	persistentVendorDatabase = vendorDB

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", handlerRoot)
	stdErr.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Server.HTTPPort), router))

	devMessage("Leaving serverMain()")
}

func handlerRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s\n", toolID)
	fmt.Fprintf(w, "%d unique OUIs and %d unique vendors in database\n", len(persistentOUIDatabase.OUIDatabase), len(persistentVendorDatabase))
}
