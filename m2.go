package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	bimg "gopkg.in/h2non/bimg.v1"
)

var (
	conf *config
)

/**
 * Serve raw image butes.
 */
func rawImageHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	if p.ByName("rawImage") == "" {
		return
	}

	fileName := conf.RawDir + "/" + p.ByName("rawImage")

	// check if file exissts
	img, err := os.Open(fileName)
	defer img.Close()

	if err != nil {
		fmt.Printf("RAW: File '%s' not found\n", fileName)
		fmt.Fprintf(w, "file not found.")
		return
	}

	// Create byte array of that can hold the image file.
	fileInfo, _ := img.Stat()
	var size = fileInfo.Size()
	bytes := make([]byte, size)

	// read file into bytes
	buffer := bufio.NewReader(img)
	_, err = buffer.Read(bytes)
	if err != nil {
		fmt.Println("RAW: Error reading buffer into byte array.")
		return
	}

	contentType := http.DetectContentType(bytes)
	// Verify that the file is an image file.
	if !strings.Contains(contentType, "image") {
		fmt.Printf("RAW: content-type = %s\n", contentType)
		fmt.Printf("RAW: The requested file '%s' does not appear to be a valid image file.\n", fileName)
		return
	}

	// Send custom headers and raw image bytes.
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.Itoa(int(size)))
	w.Header().Set("ETag", strconv.Itoa(int(time.Now().Unix()))) // Make uniquer
	w.Header().Set("Cache-Control", "max-age=2592000")           // 30 days
	w.Write(bytes)

}

/**
 * Process image based on flavor and serve the image bytes.
 */
func procImageHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	if p.ByName("procImage") == "" {
		return
	}

	re := regexp.MustCompile(`(.*)\.([^\.]*)\.([^\.]*)$`)
	parts := re.FindAllStringSubmatch(p.ByName("procImage"), 1)

	// If the requestURI does not match the pattern send error.
	if len(parts[0]) != 4 {
		w.WriteHeader(200)
		fmt.Fprintf(w, "Invalid image request\n")
		return
	}

	//  filename := parts[0][1]
	//  extension := parts[0][3]
	//  flavor := parts[0][2]

	fileName := conf.ProcDir + "/" + p.ByName("procImage")

	// check if file exissts
	img, err := os.Open(fileName)
	defer img.Close()

	// If the image doesn't exist. Is there cleaner way to test?
	if err != nil {
		// Image not found try to create it.
		// fmt.Printf("PROC: File '%s' not found\n", fileName)
		// w.WriteHeader(404)
		// fmt.Fprintf(w, "file not found.")
		// Create byte array of that can hold the image file.
		fileName = conf.RawDir + "/" + parts[0][1] + "." + parts[0][3]
		fmt.Printf("fileName: %v\n", fileName)
		img, err := os.Open(fileName)
		defer img.Close()
		if err != nil {
			fmt.Println("Error opening Raw asset.")
		}
		fileInfo, _ := img.Stat()
		var size = fileInfo.Size()
		bytes := make([]byte, size)

		// read file into bytes
		buffer := bufio.NewReader(img)
		_, err = buffer.Read(bytes)
		if err != nil {
			fmt.Println("PROC: Error reading buffer into byte array.")
			return
		}

		// Process image according to flavor.
		// newImage, err := processImage(fileName, parts[0][2])
		newImage, err := bimg.NewImage(bytes).Resize(100, 100)
		if err != nil {
			fmt.Println("Error processing file.")
			return
		}
		contentType := http.DetectContentType(bytes)
		// Verify that the file is an image file.
		if !strings.Contains(contentType, "image") {
			fmt.Printf("PROC: content-type = %s\n", contentType)
			fmt.Printf("PROC: The requested file '%s' does not appear to be a valid image file.\n", fileName)
			return
		}

		// Processed image can be served even if it isn't written to file system.
		// Send custom headers and raw image bytes.
		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Content-Length", strconv.Itoa(int(len(newImage))))
		w.Header().Set("ETag", strconv.Itoa(int(time.Now().Unix()))) // Make uniquer
		w.Header().Set("Cache-Control", "max-age=2592000")           // 30 days
		w.Write(newImage)

		ioutil.WriteFile(conf.ProcDir+"/"+p.ByName("procImage"), newImage, 0755)

	}

}
func indexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(200)
	fmt.Fprint(w, "ok\n")
}

func init() {
	conf = newConf()
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

	http.ListenAndServe("localhost:9001", router)

}
