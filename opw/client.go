package opw

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

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
	passphrase, err := c.cfg.Passphrase()
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

	re := regexp.MustCompile(`export (OP_SESSION_\w+)=\"(.*?)\"`)
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		s := scanner.Text()
		m := re.FindAllStringSubmatch(s, -1)
		if len(m) > 0 {
			logrus.Debugf("setting environment '%s' with '%s'", m[0][1], m[0][2])
			os.Setenv(m[0][1], m[0][2])
		}
	}

	return nil
}
