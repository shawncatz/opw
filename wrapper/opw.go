package wrapper

import "github.com/shawncatz/opw/config"

func NewClient(cfg *config.Config) (*Client, error) {
	return &Client{cfg}, nil
}
