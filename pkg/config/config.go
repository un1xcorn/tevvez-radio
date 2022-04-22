package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Token  string `json:"token"`
	Color  int    `json:"color"`
	Prefix string `json:"prefix"`
}

func Load(name string) (Config, error) {
	file, err := ioutil.ReadFile(name)
	if err != nil {
		return Config{}, err
	}

	c := Config{}
	err = json.Unmarshal(file, &c)
	if err != nil {
		return Config{}, err
	}

	return c, nil
}
