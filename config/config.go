package config

import (
	"encoding/json"
	"io"
	"os"

	"github.com/amirkhaki/crossword/key"
	"github.com/amirkhaki/crossword/user"
	"github.com/charmbracelet/lipgloss"
)

type Colors struct {
	QuestionBorderColor   lipgloss.Color `json:"question_border_color"`
	QuestionTextColor     lipgloss.Color `json:"question_text_color"`
	TableEditableKeyColor lipgloss.Color `json:"table_editable_key_color"`
	TableSelectedKeyColor lipgloss.Color `json:"table_selected_key_color"`
	PassPhraseKeyColor    lipgloss.Color `json:"pass_phrase_key_color"`
}

type Config struct {
	Games  []Game      `json:"games"`
	Users  []user.User `json:"users"`
	Colors Colors      `json:"colors"`
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
	InitialCol int      `json:"initial_col"`
	InitialRow int      `json:"initial_row"`
	Questions  []string `json:"questions"`
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
