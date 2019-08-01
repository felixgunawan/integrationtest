package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var MainConfig *Config

type Config struct {
	DB struct {
		Host string `json:"host"`
		Port string `json:"port"`
		User string `json:"user"`
		Pass string `json:"pass"`
		Name string `json:"name"`
	} `json:"db"`
}

func init() {
	var byteValue []byte
	if os.Getenv("EXAMPLEENV") == "" || os.Getenv("EXAMPLEENV") == "local" {
		jsonFile, err := os.Open("./config/local.json")
		if err != nil {
			log.Fatalln(err)
		}
		defer jsonFile.Close()
		byteValue, _ = ioutil.ReadAll(jsonFile)
	} else if os.Getenv("EXAMPLEENV") == "production" {
		jsonFile, err := os.Open("./config/production.json")
		if err != nil {
			log.Fatalln(err)
		}
		defer jsonFile.Close()
		byteValue, _ = ioutil.ReadAll(jsonFile)
	}
	json.Unmarshal([]byte(byteValue), &MainConfig)
}
