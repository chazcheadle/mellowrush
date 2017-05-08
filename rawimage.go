package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
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
