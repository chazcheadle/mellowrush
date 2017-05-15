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
	OrigDir  string `yaml:"origdir"`
	OrigRoot string `yaml:"origroot"`
	ProcDir  string `yaml:"procdir"`
	ProcRoot string `yaml:"procroot"`
	Defaults struct {
		Algorithm         string `yaml:"algorithm"`
		Quality           int    `yaml:"quality"`
		FallbackOoriginal bool   `yaml:"fallback-original"`
		MaxAge            int    `yaml:"max-age"`
		StoreCustom       bool   `yaml:"store-custom"`
	} `yaml:"defaults"`
	Flavors map[string]string `yaml:"flavors"`
}

func getConf() *config {
	conf := parseConfigFile()
	fmt.Printf("%v\n", conf)
	return conf
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
