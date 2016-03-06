package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Namespace string `json:"namespace,omitempty"`
	Profile   string `json:"profile"`
	Prefix    string `json:"prefix"`
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
