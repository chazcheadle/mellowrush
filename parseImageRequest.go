package main

import (
	"fmt"
	"regexp"
	"strconv"

	bimg "gopkg.in/h2non/bimg.v1"
)

func parseImageRequest(flavor string) (bimg.Options, error) {

	rawRequest := conf.Flavors[flavor]
	if rawRequest == "" {
		rawRequest = flavor
	}

	opts := bimg.Options{}

	// opts, err := parseRawTransform(rawRequest)
	fmt.Printf("PARSE: request = %v\n", flavor)
	re := regexp.MustCompile(`(\d+);(\d+);(\d+);(\d+);(\d+)`)
	parts := re.FindAllStringSubmatch(rawRequest, 1)
	if parts != nil && len(parts[0]) != 6 {
		fmt.Println("PARSE: Not a valid raw transform request.")
	} else {
		opts.Width = bimgOptConvert(parts[0][1])
		opts.Height = bimgOptConvert(parts[0][2])
		opts.Crop = true
		//	opts.Force = true
		//	opts.Extend = 1
		// opts.algorithm = bimgOptConvert(parts[0][3])
		// opts.quality = bimgOptConvert(parts[0][4])
		// opts.focalRect = bimgOptConvert(parts[0][5])
	}
	return opts, nil

}

/**
 * Utility function convert string to bimg compliant Option struct parameter.
 */
func bimgOptConvert(str string) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		val = 0
	}
	return val
}
