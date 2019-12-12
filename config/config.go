package config

import (
	"io/ioutil"
)

type Config struct {
	Email          string
	Secret         string
	Subdomain      string
	Cache          string
	PassphraseFile string `mapstructure:"passphrase_file"`
	Aliases        map[string]string
}

func (c *Config) Passphrase() (string, error) {
	data, err := ioutil.ReadFile(c.PassphraseFile)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
