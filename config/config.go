package config

import (
	"encoding/json"
	"io"
	"os"

	"github.com/amirkhaki/crossword/key"
	"github.com/charmbracelet/lipgloss"
)

type Config struct {
	Games []Game `json:"games"`
}

type Game struct {
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
	Questions             []string       `json:"questions"`
	QuestionBorderColor   lipgloss.Color `json:"question_border_color"`
	QuestionTextColor     lipgloss.Color `json:"question_text_color"`
	TableEditableKeyColor lipgloss.Color `json:"table_editable_key_color"`
	TableSelectedKeyColor lipgloss.Color `json:"table_selected_key_color"`
	PassPhraseKeyColor    lipgloss.Color `json:"pass_phrase_key_color"`
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
