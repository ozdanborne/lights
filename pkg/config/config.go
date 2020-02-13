package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

var configPath = filepath.Join(os.Getenv("HOME"), ".lights")

type config struct {
	User string
}

func (c *config) Save() error {
	return ioutil.WriteFile(configPath, []byte(c.User), 0644)
}

func Load() config {
	c := config{}
	bits, _ := ioutil.ReadFile(configPath)
	c.User = string(bits)
	return c
}
