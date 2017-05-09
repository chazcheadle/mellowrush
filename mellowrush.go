package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var (
	conf    *config
	flavors *flavorMap
)

func indexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(200)
	fmt.Fprint(w, "ok\n")
}

func init() {
	conf = getConf()
	flavors = getFlavors()
}

func main() {

	// Instantiate a new router.
	router := httprouter.New()

	// Index route.
	router.HEAD("/", indexHandler)

	// Original image route.
	router.HEAD(conf.OrigRoot+"/:origImage", origImageHandler)
	router.GET(conf.OrigRoot+"/:origImage", origImageHandler)

	// Processed image route.
	router.HEAD(conf.ProcRoot+"/:procImage", procImageHandler)
	router.GET(conf.ProcRoot+"/:procImage", procImageHandler)
	fmt.Println("start server...")
	http.ListenAndServe(conf.Hostname+":"+conf.Port, router)

}
