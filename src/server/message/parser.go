package message

import "encoding/json"

func Parse(msg []byte) (*Container, error) {
	var m Container
	err := json.Unmarshal(msg, &m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
