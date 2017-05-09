package main

import "fmt"

type options struct {
	height    int
	width     int
	algorithm int
	quality   int
	focalRect int
}

func parseImageRequest(flavor string) (*options, error) {
	fmt.Println("Parsing image options.")
	//	flavor_parts := strings.Split(flavor, ";")
	// Test for correct number of arguments
	// 	if len(flavor_parts) == 5 {
	// 		for f := range flavor_parts {
	// 			// test that each argument is a valid parameter
	// 			if !strconv.Atoi(f) {
	// 				return
	// 			}
	// 		}
	// 	}
	opts := &options{}
	opts.height = 320
	opts.width = 240
	opts.algorithm = 5
	opts.quality = 70
	opts.focalRect = 0
	return opts, nil
}
