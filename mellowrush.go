package main

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/julienschmidt/httprouter"
	//	bimg "gopkg.in/h2non/bimg.v1"
)

/*
var fileName = "hair.jpg"

func resize() {
	buffer, err := bimg.Read(fileName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	contentType := http.DetectContentType(buffer)
	if strings.Contains(contentType, "image") {
		fmt.Printf("LOG: The requested file '%s' does not appear to be a valid image file.\n", fileName)
		return
	}

	fmt.Println("This is a valid image file.")

	newImage, err := bimg.NewImage(buffer).Resize(100, 100)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	size, err := bimg.NewImage(newImage).Size()
	fmt.Println("Image height:", size.Height)
	fmt.Println("Image width:", size.Width)

	bimg.Write("hair2.jpg", newImage)
}

func resizeHandler(w http.ResponseWriter, r *http.Request) {

	resize()
}*/

type (
	ImageController struct{}
)

func NewImageController() *ImageController {
	return &ImageController{}
}

func (ic ImageController) GetRawImage(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

}

func (ic ImageController) GetProcImage(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Println(r.RequestURI)

	// Regex: (.*)\.([^\.]*)\.([^\.]*)$
	re := regexp.MustCompile(`(.*)\.([^\.]*)\.([^\.]*)$`)
	parts := re.FindAllStringSubmatch(p.ByName("procImage"), 1)
	fmt.Println(len(parts[0]))
	// If the requestURI does not match the pattern send error.
	if parts == nil || len(parts[0]) != 4 {
		// Log this.
		fmt.Fprintf(w, "Invalid image request\n")

		// Return message to browser
		w.WriteHeader(200)
		fmt.Fprintf(w, "The requested resource '%s' is invalid.\n", parts[0][0])
		return
	}

	filename := parts[0][1]
	extension := parts[0][3]
	flavor := parts[0][2]

	fmt.Println("Extension:", parts[0][3])
	fmt.Println("AIMS flavor:", parts[0][2])
	fmt.Println("File name:", parts[0][1])

	w.WriteHeader(200)
	fmt.Fprintf(w, "Filename: %s\n", filename)
	fmt.Fprintf(w, "Extension: %s\n", extension)
	fmt.Fprintf(w, "Flavor: %s\n", flavor)

	// Look for file on disk.
	// If it is there, direct the user there (ServerFile?)
	// If it is not, create it, and then direct the user.

}

func main() {

	// Instantiate a new router.
	r := httprouter.New()

	// Get a UserController instance
	ic := NewImageController()

	r.GET("/i/:rawImage", ic.GetRawImage)
	r.GET("/j/:procImage", ic.GetProcImage)

	http.ListenAndServe("localhost:9001", r)

}
