package main

import (
	"errors"
	"log"

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

func toUpper(r rune) (rune, error)  {
	if r >= 'A' && r <= 'Z' {
		return r, nil
	}

	if r >= 'a' && r <= 'z' {
		return r - 'a' + 'A', nil
	}

	return r, errors.New("the rune is not an alphabet")

}

type game struct {
	actual key.Table
	mustBe key.Table
	crrntRow int
	crrntCol int
}

func (g *game) Init() tea.Cmd {
	return nil
}

func (g *game) View() string {
	var rows [key.COLS]string
	for i := 0; i < key.ROWS; i++ {
		var cols [key.COLS]string
		for j := 0; j < key.COLS; j++ {
			if i == g.crrntRow && j == g.crrntCol {
				cols[j] = g.actual[i][j].Render(colorPrimary)
			} else {
				cols[j] = g.actual[i][j].Render(colorSecondary)
			}
		}
		rows[i] = lipgloss.JoinHorizontal(lipgloss.Bottom, cols[:]...)
	}
	return lipgloss.JoinVertical(lipgloss.Center, rows[:]...)
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
	if g.crrntRow + 1 == key.ROWS {
		return nil
	}
	if g.actual[g.crrntRow+1][g.crrntCol] == key.LOCKED {
		return nil
	}
	g.crrntRow++

	return nil
}

func (g *game) goUp() tea.Cmd {
	if g.crrntRow == 0 {
		return nil
	}
	if g.actual[g.crrntRow-1][g.crrntCol] == key.LOCKED {
		return nil
	}
	g.crrntRow--

	return nil
}

func (g *game) goLeft() tea.Cmd {
	if g.crrntCol == 0 {
		return nil
	}
	if g.actual[g.crrntRow][g.crrntCol-1] == key.LOCKED {
		return nil
	}
	g.crrntCol--

	return nil
}

func (g *game) goRight() tea.Cmd {
	if g.crrntCol + 1 == key.COLS {
		return nil
	}
	if g.actual[g.crrntRow][g.crrntCol+1] == key.LOCKED {
		return nil
	}
	g.crrntCol++

	return nil
}

func (g *game) insertKey(r rune) tea.Cmd {
	r, err := toUpper(r)
	if err != nil {
		return nil
	}
	k := key.Letters[r]
	g.actual[g.crrntRow][g.crrntCol] = k
	return nil	
}

func (g *game) Ended() bool {
	for i:=0; i < key.COLS; i++ {
		for j := 0; j < key.ROWS; j++ {
			if g.actual[i][j] != g.mustBe[i][j] {
				return false
			}
		}
	}
	return true	
}
// TODO read initial position from config (crrntRow, crrntCol) which must not be a locked key
// TODO reas mustBe and actual from config 
func NewGame(cfg config.Config) ( *game, error ) {
	g := game{}
	for i:=0; i < key.COLS; i++ {
		for j := 0; j < key.ROWS; j++ {
			g.actual[i][j] = key.EMPTY
		}
	}
	g.actual[1][1] = key.LOCKED
	g.actual[5][10] = key.LOCKED
	g.actual[12][12] = key.Z
	return &g, nil
}

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	g, err := NewGame(cfg)
	if err != nil {
		log.Fatal(err)
	}
	p := tea.NewProgram(g)
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
