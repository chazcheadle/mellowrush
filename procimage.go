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
	fmt.Printf("PROC: requeseted filename: %v\n", fileName)
	// check if file exissts
	img, err := os.Open(fileName)
	defer img.Close()

	// If the image doesn't exist. Is there cleaner way to test?
	if err != nil {
		// Image not found try to create it.
		// fmt.Printf("PROC: File '%s' not found\n", fileName)
		// w.WriteHeader(404)
		// fmt.Fprintf(w, "file not found.")

		// Try to open the original image.
		fileName = conf.OrigDir + "/" + parts[0][1] + "." + parts[0][3]
		origImg, err := os.Open(fileName)
		defer origImg.Close()
		if err != nil {
			fmt.Printf("Error opening original image asset: %s.\n", fileName)
		}

		// Create byte array of that can hold the image file.
		fileInfo, _ := origImg.Stat()
		var size = fileInfo.Size()
		bytes := make([]byte, size)

		// read file into bytes
		buffer := bufio.NewReader(origImg)
		_, err = buffer.Read(bytes)
		if err != nil {
			fmt.Println("PROC: Error reading buffer into byte array.")
			return
		}

		// Process image according to flavor.
		// newImage, err := processImage(fileName, parts[0][2])
		options, err := parseImageRequest(parts[0][2])
		if err != nil {
			fmt.Println("Error parsing image processing parameters.")
			return
		}

		// Send resize request
		// replace with bimgOpts?
		newImage, err := bimg.NewImage(bytes).ResizeAndCrop(options.Height, options.Width)
		//newImage, err := bimg.NewImage(bytes).Process(options)
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
		// Send custom headers and image bytes.
		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Content-Length", strconv.Itoa(int(len(newImage))))
		w.Header().Set("ETag", strconv.Itoa(int(time.Now().Unix()))) // Make uniquer
		w.Header().Set("Cache-Control", "max-age=2592000")           // 30 days
		w.Write(newImage)

		// Check if the image should be written to disk.
		if conf.Defaults.StoreCustom {
			ioutil.WriteFile(conf.ProcDir+"/"+p.ByName("procImage"), newImage, 0755)
		}
	} else {
		// Serve file
		fmt.Printf("Serving existing asset: %s\n", fileName)
		if err != nil {
			fmt.Printf("Error opening processed image asset: %s.\n", fileName)
			return
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

		contentType := http.DetectContentType(bytes)
		// Verify that the file is an image file.
		if !strings.Contains(contentType, "image") {
			fmt.Printf("PROC: content-type = %s\n", contentType)
			fmt.Printf("PROC: The requested file '%s' does not appear to be a valid image file.\n", fileName)
			return
		}

		// Processed image can be served even if it isn't written to file system.
		// Send custom headers and image bytes.
		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Content-Length", strconv.Itoa(int(size)))
		w.Header().Set("ETag", strconv.Itoa(int(time.Now().Unix())))                   // Make uniquer
		w.Header().Set("Cache-Control", "max-age="+strconv.Itoa(conf.Defaults.MaxAge)) // 30 days
		w.Write(bytes)
	}

}
