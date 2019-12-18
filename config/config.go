package config

import (
	"errors"
	"io/ioutil"
	"strings"

	"github.com/sirupsen/logrus"

	keyring "github.com/zalando/go-keyring"
)

type Config struct {
	Email      string
	Secret     string
	Subdomain  string
	Cache      string
	Passphrase string
	Aliases    map[string]string
	Debug      bool
}

func (c *Config) GetSecret() (string, error) {
	return c.getValue(c.Secret)
}

func (c *Config) GetPassphrase() (string, error) {
	return c.getValue(c.Passphrase)
}

func (c *Config) getValue(value string) (string, error) {
	v := strings.Split(value, ":")
	logrus.Debugf("value setting: %#v", v)
	switch v[0] {
	case "file":
		return c.getValueFile(v[1])
	case "keychain":
		return c.getValueKeychain(v[1] + ":" + v[2])
	}
	return "", errors.New("unknown value type: " + v[0])
}

func (c *Config) getValueFile(file string) (string, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (c *Config) getValueKeychain(setting string) (string, error) {
	v := strings.Split(setting, ":")
	secret, err := keyring.Get(v[0], v[1])
	if err != nil {
		return "", err
	}

	return secret, nil
}
