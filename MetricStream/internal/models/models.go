package models

import (
	"encoding/json"
)

type Data struct {
	Cpu  CpuInfo
	Vram VramInfo
}

func (c *Data) Marshal() ([]byte, error) {
	jsonData, err := json.Marshal(c)
	if err != nil {
		return []byte(""), err
	}
	return []byte(jsonData), nil
}
func (c *Data) Unmarshal(b []byte) {
	err := json.Unmarshal(b, c)
	if err != nil {
		return
	}
}

type Settings struct {
	UpdateTime int
}
