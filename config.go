package main

import (
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type config struct {
	Hostname string `yaml:"hostname"`
	Port     string `yaml:"port"`
	OrigDir  string `yaml:"origDir"`
	ProcDir  string `yaml:"procDir"`
}

type flavorMap struct {
	Flavors map[string]string
}

func getConf() *config {
	conf := parseConfigFile()
	return conf
}

func getFlavors() *flavorMap {
	flavorConf := parseFlavorsFile()
	fmt.Printf("%v\n", flavorConf)
	return flavorConf
}

func parseConfigFile() *config {

	confFile := "mellowrush.yml"
	if len(os.Args[1:]) > 0 {
		confFile = os.Args[1]
	}
	conf := &config{}

	f, err := os.Open(confFile)
	defer f.Close()
	if err != nil {
		//Log error.
		fmt.Println("Couldn't open config file.")
		return conf
	}
	d, err := ioutil.ReadAll(f)
	if err != nil {
		// Log error.
		fmt.Println("Couldn't read config file.")
		return conf
	}
	err = yaml.Unmarshal(d, conf)
	if err != nil {
		// Log error.
		fmt.Println("Couldn't unmarshal config file yaml.")
		return conf
	}
	return conf
}

func parseFlavorsFile() *flavorMap {

	flavorFile := "flavors.yml"
	if len(os.Args[1:]) > 0 {
		flavorFile = os.Args[1]
	}
	flavors := &flavorMap{}

	f, err := os.Open(flavorFile)
	defer f.Close()
	if err != nil {
		//Log error.
		fmt.Println("Couldn't open flavor config file.")
		return flavors
	}
	d, err := ioutil.ReadAll(f)
	if err != nil {
		// Log error.
		fmt.Println("Couldn't read flavor config file.")
		return flavors
	}
	err = yaml.Unmarshal(d, flavors)
	if err != nil {
		// Log error.
		fmt.Println("Couldn't unmarshal flavor config file yaml.")
		return flavors
	}
	return flavors

}
