package main

import (
	"gopkg.in/yaml.v2"
	"os"
)

/*
	Global config object
*/
var (
	global_cfg *AvologoConfig
	global_cfg_path string
)

/*
	Struct defining an avologo configuration (read from file)
*/
type AvologoConfig struct {
    Server struct {
		Host 		string `yaml:"host"`
        Port 		int `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Host 		string `yaml:"host"`
		Port 		int `yaml:"port"`
		User 		string `yaml:"user"`
		Password 	string `yaml:"password"`
		DBName 		string `yaml:"dbname"`
	} `yaml:"database"`
	Client struct {
		Destination		string `yaml:"destination"`
		FriendlyName	string `yaml:"friendly_name"`
		Watch 			[]string `yaml:"watch"`
		ErrorKeywords 	[]string `yaml:"error_keywords"`
		WarningKeywords []string `yaml:"warning_keywords"`
    } `yaml:"client"`
}

/*
	Parse config file into AvologoConfig struct
*/
func parseConfig(path string) *AvologoConfig {
	f, err := os.Open(path)
	if (err != nil) {
		panic(err)
	}
	defer f.Close()
	
	cfg := new(AvologoConfig)
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if (err != nil) {
		panic(err)
	}
	return cfg
}