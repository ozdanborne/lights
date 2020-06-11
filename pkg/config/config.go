package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

var configPath = filepath.Join(os.Getenv("HOME"), ".lights")

type config struct {
	User  string
	Group int
}

func (c *config) Save() error {
	bits, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(configPath, bits, 0644)
}

func Load() (config, error) {
	c := config{}

	bits, err := ioutil.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return c, nil
		}
		return c, err
	}
	if err := json.Unmarshal(bits, &c); err != nil {
		return c, err
	}

	return c, nil
}
