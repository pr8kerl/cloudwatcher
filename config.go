package main

import (
	"encoding/json"
	"io/ioutil"
)

type Account struct {
	AccountId   int    `json:"accountId"`
	AccountName string `json:"accountName"`
	AccessKey   string `json:"accessKey"`
	SecretKey   string `json:"secretKey"`
}

type GraphiteCfg struct {
	InstanceStem string `json:"instancestem"`
}

type Config struct {
	Accounts []Account   `json:"accounts"`
	Regions  []string    `json:"regions"`
	Debug    bool        `json:"debug,omitempty"`
	Graphite GraphiteCfg `json:"graphite,omitempty"`
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
