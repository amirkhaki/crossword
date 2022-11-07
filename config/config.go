package config

import (
	"encoding/json"
	"io"
	"os"

	"github.com/amirkhaki/crossword/key"
)

type Config struct {
	Rows   int `json:"rows"`
	Cols   int `json:"cols"`
	Actual struct {
		Keys []struct {
			Row int     `json:"row"`
			Col int     `json:"col"`
			Key key.Key `json:"key"`
		} `json:"keys"`
	} `json:"actual"`
	InitialCol int `json:"initial_col"`
	InitialRow int `json:"initial_row"`
	Mustbe     struct {
		Keys []struct {
			Row int     `json:"row"`
			Col int     `json:"col"`
			Key key.Key `json:"key"`
		} `json:"keys"`
	} `json:"mustbe"`
}

func New(path string) (cfg Config, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	data, err := io.ReadAll(f)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &cfg)
	return
}
