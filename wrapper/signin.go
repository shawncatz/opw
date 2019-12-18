package wrapper

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/shawncatz/opw/simplecache"
)

// SignInFull does a full signin using all arguments and forces the cache to reset
func (c *Client) SignInFull() error {
	return c.signIn(true)
}

// SignIn does a simple signin using just the subdomain
func (c *Client) SignIn() error {
	return c.signIn(false)
}

// signIn is the utility function for the main SignIn* functions
func (c *Client) signIn(full bool) error {
	value, err := c.getSession(full)
	if err != nil {
		return err
	}

	m := strings.Split(value, ":")
	logrus.Debugf("setting %s to %s", m[0], m[1])
	if err := os.Setenv(m[0], m[1]); err != nil {
		return err
	}

	return nil
}

// getSession fetches the session from the cache, optionally forcing a reset with full = true
func (c *Client) getSession(full bool) (string, error) {
	cache, err := simplecache.New(c.cfg.Cache)
	if err != nil {
		return "", err
	}

	return cache.Fetch("session", 30*time.Minute, full, func() (string, error) {
		passphrase, err := c.cfg.GetPassphrase()
		if err != nil {
			return "", err
		}

		secret, err := c.cfg.GetSecret()
		if err != nil {
			return "", err
		}

		args := []string{"signin"}
		if full {
			args = append(args, c.cfg.Subdomain+".1password.com", c.cfg.Email, secret)
		} else {
			args = append(args, c.cfg.Subdomain)
		}

		cmd := exec.Command("/usr/local/bin/op", args...)
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return "", err
		}

		go func() {
			defer stdin.Close()
			io.WriteString(stdin, passphrase)
		}()

		out, err := cmd.CombinedOutput()
		if err != nil {
			return "", err
		}

		re := regexp.MustCompile(`export (OP_SESSION_\w+)=\"(.*?)\"`)
		scanner := bufio.NewScanner(strings.NewReader(string(out)))
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
