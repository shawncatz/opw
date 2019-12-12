package opw

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/shawncatz/opw/simplecache"

	"github.com/sirupsen/logrus"

	"github.com/shawncatz/opw/config"
)

type Client struct {
	cfg *config.Config
}

func (c *Client) GetItem(uuid string) (*Item, error) {
	out, err := exec.Command("/usr/local/bin/op", "get", "item", uuid).Output()
	if err != nil {
		return nil, err
	}

	item := &Item{}
	err = json.Unmarshal([]byte(out), item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (c *Client) SignIn() error {
	passphrase, err := c.cfg.GetPassphrase()
	if err != nil {
		return err
	}

	cmd := exec.Command("/usr/local/bin/op", "signin", c.cfg.Subdomain)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, passphrase)
	}()

	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	value, err := c.getSession(string(out))
	if err != nil {
		return err
	}

	m := strings.Split(value, ":")
	if err := os.Setenv(m[0], m[1]); err != nil {
		return err
	}

	return nil
}

func (c *Client) getSession(out string) (string, error) {
	cache, err := simplecache.New(c.cfg.Cache)
	if err != nil {
		return "", err
	}

	return cache.Fetch("session", 30*time.Minute, func() (string, error) {
		re := regexp.MustCompile(`export (OP_SESSION_\w+)=\"(.*?)\"`)
		scanner := bufio.NewScanner(strings.NewReader(out))
		for scanner.Scan() {
			s := scanner.Text()
			m := re.FindAllStringSubmatch(s, -1)
			if len(m) > 0 {
				logrus.Debugf("setting environment '%s' with '%s'", m[0][1], m[0][2])
				return m[0][1] + ":" + m[0][2], nil
			}
		}
		return "", nil
	})
}
