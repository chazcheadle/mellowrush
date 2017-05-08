package main

type config struct {
	RawDir  string `yaml:"rawDir"`
	ProcDir string `yaml:"prodDir"`
}

func newConf() *config {
	conf := parseConfigFile()
	return conf
}

func parseConfigFile() *config {
	conf := &config{}
	conf.RawDir = "/tmp/i"
	conf.ProcDir = "/tmp/j"
	return conf
}
