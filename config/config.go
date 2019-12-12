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
	Passphrase string `mapstructure:"passphrase"`
	Aliases    map[string]string
}

func (c *Config) GetPassphrase() (string, error) {
	v := strings.Split(c.Passphrase, ":")
	logrus.Infof("passphrase setting: %#v", v)
	switch v[0] {
	case "file":
		return c.getPassphraseFile(v[1])
	case "keychain":
		return c.getPassphraseKeychain(v[1] + ":" + v[2])
	}

	return "", errors.New("unknown passphrase type: " + v[0])
}

func (c *Config) getPassphraseFile(file string) (string, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (c *Config) getPassphraseKeychain(setting string) (string, error) {
	v := strings.Split(setting, ":")
	secret, err := keyring.Get(v[0], v[1])
	if err != nil {
		return "", err
	}

	return secret, nil
}
