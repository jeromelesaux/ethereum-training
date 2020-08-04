package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

//
// ethereum endpoint defined here
// https://infura.io/dashboard/ethereum/7de903803c31428bbdd1186107a2d660/settings
//
type Configuration struct {
	DirectorySavePath string `json:"directorysavepath"`
	EthereumEndpoint  string `json:"ethereumendpoint"`
	PrivateKey        string `json:"privatekey"`
}

var (
	doOnce   sync.Once
	MyConfig *Configuration
)

func LoadConfig() {
	doOnce.Do(
		func() {
			MyConfig = &Configuration{
				Rootpath: "./",
			}

		})
}

func LoadConfigFile(file string) {
	doOnce.Do(
		func() {
			f, err := os.Open(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Cannot open configuration file [%s] with error :%v\n", file, err)
				log.Fatal()
			}
			defer f.Close()
			c := &Configuration{}
			if err := json.NewDecoder(f).Decode(c); err != nil {
				fmt.Fprintf(os.Stderr, "Cannot decode configuration file [%s] with error :%v\n", file, err)
				log.Fatal()
			}
			MyConfig = c
			fmt.Fprintf(os.Stdout, "Loaded configuration from file [%s]\n", file)
		})
	return
}

func (c *Configuration) GetFilepaths() string {
	f := filepath.Join(c.DirectorySavePath, "tmpfiles")
	os.MkdirAll(f, os.ModePerm)
	return f
}

func (c *Configuration) GetHashpaths() string {
	f := filepath.Join(c.DirectorySavePath, "hashes")
	os.MkdirAll(f, os.ModePerm)
	return f
}
