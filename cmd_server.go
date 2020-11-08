package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

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
	router.HandleFunc("/mac/{id}", handlerMAC)
	router.HandleFunc("/vendor/{id}", handlerVendor)
	stdErr.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Server.HTTPPort), router))

	devMessage("Leaving serverMain()")
}

func handlerRoot(w http.ResponseWriter, r *http.Request) {
	devMessage("Entering handlerRoot()")

	fmt.Fprintf(w, "%s\n", toolID)
	fmt.Fprintf(w, "%d unique OUIs and %d unique vendors in database\n", len(persistentOUIDatabase.OUIDatabase), len(persistentVendorDatabase))
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "Usable endpoints:\n")
	fmt.Fprintf(w, "  /mac/{id}\n")
	fmt.Fprintf(w, "  /vendor/{id}\n")

	devMessage("Leaving handlerRoot()")
}

func handlerMAC(w http.ResponseWriter, r *http.Request) {
	devMessage("Entering handlerMAC()")

	var vendorName string

	vars := mux.Vars(r)
	mac := vars["id"]

	if !isValidMAC(mac) {
		fmt.Fprintf(w, "Warning: MAC %s is invalid.\n", mac)
		return
	}
	mac = strings.ToLower(mac)
	oui, ouiErr := extractOUI(mac)
	if ouiErr != nil {
		fmt.Fprintf(w, "Warning: Unable to map OUI for MAC %s: %s\n", mac, ouiErr)
		return
	}
	if vendorName = persistentOUIDatabase.OUIDatabase[oui].VendorName; vendorName == "" {
		vendorName = "(unregistered)"
	}
	normalized, normalizedErr := normalizeMAC(mac)
	if normalizedErr != nil {
		fmt.Fprintf(w, "Warning: MAC could not be normalized: %s\n", normalizedErr)
		return
	}
	fmt.Fprintf(w, "%s = %s\n", normalized, vendorName)

	devMessage("Leaving handlerMAC()")
}

func handlerVendor(w http.ResponseWriter, r *http.Request) {
	devMessage("Entering handlerVendor()")

	vars := mux.Vars(r)
	vendor := vars["id"]
	fmt.Fprintf(w, "Would look up vendor %s now.\n", vars["id"])

	if ouis, vendorExists := persistentVendorDatabase[vendor]; vendorExists {
		for _, oui := range ouis {
			mac, macErr := normalizeMAC(oui)
			if macErr != nil {
				fmt.Fprintf(w, "Warning: MAC could not be normalized: %s\n", macErr)
				continue
			}
			fmt.Fprintf(w, "%s = %s\n", vendor, mac[:8])
		}
	}

	devMessage("Leaving handlerVendor()")
}
