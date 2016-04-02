package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Profile    string            `json:"profile"`
	Prefix     string            `json:"prefix"`
	PollPeriod int64             `json:"pollperiod"`
	Namespaces map[string]string `json:"namespaces"`
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
