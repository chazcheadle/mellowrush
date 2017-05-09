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
	fmt.Printf("%v\n", conf)
	flavors = getFlavors()
	fmt.Println("init")
}

func main() {

	// Instantiate a new router.
	router := httprouter.New()

	// Index route.
	router.HEAD("/", indexHandler)

	// Raw image route.
	router.HEAD("/i/:rawImage", rawImageHandler)
	router.GET("/i/:rawImage", rawImageHandler)

	// Processed image route.
	router.HEAD("/j/:procImage", procImageHandler)
	router.GET("/j/:procImage", procImageHandler)
	fmt.Println("start...")
	http.ListenAndServe(conf.Hostname+":"+conf.Port, router)

}
