package model

import (
	"unicode"

	"github.com/amirkhaki/crossword/config"
	"github.com/amirkhaki/crossword/key"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	colorPrimary   = lipgloss.Color("#d7dadc")
	colorSecondary = lipgloss.Color("#626262")
	colorSeparator = lipgloss.Color("#9c9c9c")
	colorYellow    = lipgloss.Color("#b59f3b")
	colorGreen     = lipgloss.Color("#538d4e")
)

type game struct {
	rows     int
	cols     int
	actual   [][]key.Key
	mustBe   [][]key.Key
	crrntRow int
	crrntCol int
}

func (g *game) Init() tea.Cmd {
	return nil
}

func (g *game) View() string {
	var rows []string = make([]string, g.rows)
	for i := 0; i < g.rows; i++ {
		var cols []string = make([]string, g.cols)
		for j := 0; j < g.cols; j++ {
			if i == g.crrntRow && j == g.crrntCol {
				cols[j] = g.actual[i][j].Render(colorPrimary)
			} else {
				cols[j] = g.actual[i][j].Render(colorSecondary)
			}
		}
		rows[i] = lipgloss.JoinHorizontal(lipgloss.Bottom, cols...)
	}
	return lipgloss.JoinVertical(lipgloss.Center, rows...)
}

func (g *game) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyRight:
			return g, g.goRight()
		case tea.KeyLeft:
			return g, g.goLeft()
		case tea.KeyUp:
			return g, g.goUp()
		case tea.KeyDown:
			return g, g.goDown()
		case tea.KeyCtrlC:
			return g, tea.Quit
		case tea.KeyRunes:
			if len(msg.Runes) == 1 {
				return g, g.insertKey(msg.Runes[0])
			}
		}
	}
	return g, nil
}

func (g *game) goDown() tea.Cmd {
	if g.crrntRow+1 == g.rows {
		return nil
	}
	if g.actual[g.crrntRow+1][g.crrntCol].State == key.READONLY {
		return nil
	}
	g.crrntRow++

	return nil
}

func (g *game) goUp() tea.Cmd {
	if g.crrntRow == 0 {
		return nil
	}
	if g.actual[g.crrntRow-1][g.crrntCol].State == key.READONLY {
		return nil
	}
	g.crrntRow--

	return nil
}

func (g *game) goLeft() tea.Cmd {
	if g.crrntCol == 0 {
		return nil
	}
	if g.actual[g.crrntRow][g.crrntCol-1].State == key.READONLY {
		return nil
	}
	g.crrntCol--

	return nil
}

func (g *game) goRight() tea.Cmd {
	if g.crrntCol+1 == g.cols {
		return nil
	}
	if g.actual[g.crrntRow][g.crrntCol+1].State == key.READONLY {
		return nil
	}
	g.crrntCol++

	return nil
}

func (g *game) insertKey(r rune) tea.Cmd {
	if !unicode.IsLetter(r) || g.actual[g.crrntRow][g.crrntCol].State == key.READONLY {
		return nil
	}
	r = unicode.ToUpper(r)
	k := key.Letters[r]
	g.actual[g.crrntRow][g.crrntCol].Char = k
	if g.Ended() {
		return tea.Quit
	}
	return nil
}

func (g *game) Ended() bool {
	for i := 0; i < g.rows; i++ {
		for j := 0; j < g.cols; j++ {
			if !g.mustBe[i][j].IsEqual(g.actual[i][j]) {
				return false
			}
		}
	}
	return true
}

func NewGame(cfg config.Config) (*game, error) {
	g := game{}
	g.rows = cfg.Rows
	g.cols = cfg.Cols
	g.actual = make([][]key.Key, cfg.Rows)
	for i := 0; i < cfg.Rows; i++ {
		g.actual[i] = make([]key.Key, cfg.Cols)
	}
	g.mustBe = make([][]key.Key, cfg.Rows)
	for i := 0; i < cfg.Rows; i++ {
		g.mustBe[i] = make([]key.Key, cfg.Cols)
	}
	for _, v := range cfg.Actual.Keys {
		g.actual[v.Row][v.Col] = v.Key
	}
	for _, v := range cfg.Mustbe.Keys {
		g.mustBe[v.Row][v.Col] = v.Key
	}
	g.crrntCol = cfg.InitialCol
	g.crrntRow = cfg.InitialRow
	return &g, nil
}

// TODO show time of other users realtime
// TODO cause everybody can register multiple times with different names and cheat \
// it would be nice if 1. start game at fixed time or 2. have multiple with same level of difficulity
// or limit registration, so only verified users will be able to play
