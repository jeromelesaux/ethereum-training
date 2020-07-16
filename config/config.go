package config

import (
	"os"
	"path/filepath"
	"sync"
)

//
// ethereum endpoint defined here
// https://infura.io/dashboard/ethereum/7de903803c31428bbdd1186107a2d660/settings
//
type Configuration struct {
	Rootpath         string
	EthereumEndpoint string
	PrivateKey       string
}

var (
	doOnce   sync.Once
	MyConfig *Configuration
)

func LoadConfig() {
	doOnce.Do(
		func() {
			MyConfig = &Configuration{
				Rootpath:         "./",
				EthereumEndpoint: "https://ropsten.infura.io/v3/7de903803c31428bbdd1186107a2d660",
				PrivateKey:       "48218b47d9afba13df85e4b29e4e0bb73ae526cdebb316738832be607e7c7174",
			}

		})
}

func (c *Configuration) GetFilepaths() string {
	f := filepath.Join(c.Rootpath, "tmpfiles")
	os.MkdirAll(f, os.ModePerm)
	return f
}

func (c *Configuration) GetHashpaths() string {
	f := filepath.Join(c.Rootpath, "hashes")
	os.MkdirAll(f, os.ModePerm)
	return f
}
