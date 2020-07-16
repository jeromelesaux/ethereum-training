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
				Rootpath: "./",
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
