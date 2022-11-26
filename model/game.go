package model

import (
	"errors"
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
type endGameMsg struct{}
var games []*game
var currentGame int
type endScreen struct {
  
}

func (_ endScreen) Init() tea.Cmd {
  return nil
}

func (e endScreen) Update(tea.Msg) (tea.Model, tea.Cmd) {
  return e, tea.Quit
}

func (_ endScreen) View() string {
  return "game ended, do something to quit"
}

type betweenGame struct {
	updateCounter int
}

func (bg betweenGame) Init() tea.Cmd {
	return nil
}

func (bg betweenGame) Update(tea.Msg) (tea.Model, tea.Cmd) {
	bg.updateCounter++
	if bg.updateCounter > 1 {
    if currentGame >= len(games)-1 {
      return endScreen{}, nil
    }
		currentGame++
		return games[currentGame], nil
	}
  return bg, nil
}

func (bg betweenGame) View() string {
  g := games[currentGame]
	var rows []string = make([]string, g.rows)
	for i := 0; i < g.rows; i++ {
		var cols []string = make([]string, g.cols)
		for j := 0; j < g.cols; j++ {
			if g.mustBe[i][j].State == key.PASSPHRASE {
				cols[j] = g.mustBe[i][j].Render(g.passPhraseKeyColor)
			} else {
				cols[j] = g.mustBe[i][j].Render(g.keyColor)
			}
		}
		rows[i] = lipgloss.JoinHorizontal(lipgloss.Bottom, cols...)
	}
	table := lipgloss.JoinVertical(lipgloss.Center, rows...)
	notificatoins := lipgloss.JoinVertical(lipgloss.Left, "Bordered keys are passphrase letters, you'll need them later", "Press any key to continue to next game")
	notificatoins = lipgloss.NewStyle().
		Padding(0, 1).
		Border(lipgloss.NormalBorder()).
		BorderForeground(g.questionTextColor).
		Foreground(g.questionBorderColor).
		Render(notificatoins)
	board := lipgloss.JoinVertical(lipgloss.Center, table, notificatoins)
	return lipgloss.Place(g.width, g.height, lipgloss.Center, lipgloss.Center, board)
  
}

type game struct {
	width               int
	height              int
	rows                int
	cols                int
	actual              [][]key.Key
	mustBe              [][]key.Key
	crrntRow            int
	crrntCol            int
	questions           []string
	questionBorderColor lipgloss.Color
	questionTextColor   lipgloss.Color
	keyColor            lipgloss.Color
	currentKeyColor     lipgloss.Color
  passPhraseKeyColor  lipgloss.Color
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
				cols[j] = g.actual[i][j].Render(g.currentKeyColor)
			} else {
				cols[j] = g.actual[i][j].Render(g.keyColor)
			}
		}
		rows[i] = lipgloss.JoinHorizontal(lipgloss.Bottom, cols...)
	}
	table := lipgloss.JoinVertical(lipgloss.Center, rows...)
	questions := lipgloss.JoinVertical(lipgloss.Left, g.questions...)
	questions = lipgloss.NewStyle().
		Padding(0, 1).
		Border(lipgloss.NormalBorder()).
		BorderForeground(g.questionTextColor).
		Foreground(g.questionBorderColor).
		Render(questions)
	board := lipgloss.JoinHorizontal(lipgloss.Center, table, questions)
	return lipgloss.Place(g.width, g.height, lipgloss.Center, lipgloss.Center, board)
}

func (g *game) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
  case endGameMsg:
  return betweenGame{}.Update(nil)
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
	case tea.WindowSizeMsg:
		return g, g.doResize(msg)
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
		return g.EndGame()
	}
	return nil
}
func (g *game) doResize(msg tea.WindowSizeMsg) tea.Cmd {
	g.height = msg.Height
	g.width = msg.Width
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


func (g *game) EndGame() tea.Cmd {
	return func() tea.Msg {
    return endGameMsg{}
  }

}

func newGame(cfg config.Game, height, width int) (*game, error) {
	g := game{}
  g.height = height
  g.width = width
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
	g.questions = cfg.Questions
	g.questionBorderColor = cfg.QuestionBorderColor
	g.questionTextColor = cfg.QuestionTextColor
	g.currentKeyColor = cfg.TableSelectedKeyColor
	g.keyColor = cfg.TableEditableKeyColor
  g.passPhraseKeyColor = cfg.PassPhraseKeyColor
	return &g, nil
}

func NewGame(cfg config.Config, height, width int) (*game, error) {
  games = []*game{}
	for _, g := range cfg.Games {
		gm, err := newGame(g, height, width)
		if err != nil {
			return nil, err
		}
		games = append(games, gm)
	}
  if len(games) == 0 {
    return nil, errors.New("no game found in config")
  }
	currentGame = 0
	return games[0], nil
}

// TODO show time of other users realtime
// TODO cause everybody can register multiple times with different names and cheat \
// it would be nice if 1. start game at fixed time or 2. have multiple with same level of difficulity
// or limit registration, so only verified users will be able to play
