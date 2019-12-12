package opw

import (
	"encoding/json"
	"os/exec"
)

func GetItem(uuid string) (*Item, error) {
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
