package wrapper

import (
	"encoding/json"
	"os/exec"

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

func (c *Client) List() ([]*Item, error) {
	out, err := exec.Command("/usr/local/bin/op", "list", "items").Output()
	if err != nil {
		return nil, err
	}

	var items []*Item
	if err := json.Unmarshal(out, &items); err != nil {
		return nil, err
	}

	return items, nil
}
