package main

import (
	"fmt"
	"net/http"

	mux "github.com/gorilla/mux"
)

func serverMain() {
	devMessage("Entering serverMain()")
	sanitizeArguments()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", handlerRoot)
	stdErr.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Server.HTTPPort), router))

	devMessage("Leaving serverMain()")
}

func handlerRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", toolID)
}
