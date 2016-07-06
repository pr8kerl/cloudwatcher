package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Profile                  string            `json:"profile"`
	Region                   string            `json:"region"`
	Prefix                   string            `json:"prefix"`
	PollInterval             string            `json:"pollInterval"`
	Namespaces               map[string]string `json:"namespaces"`
	Debug                    bool              `json:"debug"`
	AvailableMetricsInterval int64             `json:"refreshAvailableMetricsInterval"`
}

func InitialiseConfig(cfg string) (err error) {

	// read in json file
	dat, err := ioutil.ReadFile(cfg)
	if err != nil {
		return err
	}

	// convert json to a node struct
	err = json.Unmarshal(dat, &config)
	if err != nil {
		return err
	}

	return nil
}
