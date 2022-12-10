package model

import (
	"unicode"

	"github.com/amirkhaki/crossword/config"
	"github.com/amirkhaki/crossword/data"
	"github.com/amirkhaki/crossword/key"
	"github.com/amirkhaki/crossword/user"

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

type game struct {
	err                 error
	updateCounter       int
	usr                 user.User
	width               int
	height              int
	crrntRow            int
	crrntCol            int
	questionBorderColor lipgloss.Color
	questionTextColor   lipgloss.Color
	keyColor            lipgloss.Color
	currentKeyColor     lipgloss.Color
	passPhraseKeyColor  lipgloss.Color
}

func (g *game) Init() tea.Cmd {
	return nil
}
func (g *game) afterGameView() string {
	rowCount, err := data.GetGroupRows(g.usr.Group)
	if err != nil {
		g.err = err
		return "an error accured: " + err.Error() + " press any keyboard key to exit"
	}
	colCount, err := data.GetGroupCols(g.usr.Group)
	if err != nil {
		g.err = err
		return "an error accured: " + err.Error() + " press any keyboard key to exit"
	}
	var rows []string = make([]string, rowCount)
	for i := 0; i < rowCount; i++ {
		var cols []string = make([]string, colCount)
		for j := 0; j < colCount; j++ {
			k, err := data.GetGroupRowColumn(g.usr.Group, i, j)
			if err != nil {
				g.err = err
				return "an error accured: " + err.Error() + " press any keyboard key to exit"
			}

			if k.State == key.PASSPHRASE {
				cols[j] = k.MustRender(g.passPhraseKeyColor)
			} else {
				cols[j] = k.MustRender(g.keyColor)
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

func (g *game) View() string {
	if g.err != nil {
		return "an error accured: " + g.err.Error() + " press any keyboard key to exit"
	}
	isAfterGame, err := data.IsAfterGame(g.usr.Group)
	if err != nil {
		g.err = err
		return "an error accured: " + err.Error() + " press any keyboard key to exit"
	}
	if isAfterGame {
		return g.afterGameView()
	}
	rowCount, err := data.GetGroupRows(g.usr.Group)
	if err != nil {
		g.err = err
		return "an error accured: " + err.Error() + " press any keyboard key to exit"
	}
	colCount, err := data.GetGroupCols(g.usr.Group)
	if err != nil {
		g.err = err
		return "an error accured: " + err.Error() + " press any keyboard key to exit"
	}
	var rows []string = make([]string, rowCount)
	for i := 0; i < rowCount; i++ {
		var cols []string = make([]string, colCount)
		for j := 0; j < colCount; j++ {
			k, err := data.GetGroupRowColumn(g.usr.Group, i, j)
			if err != nil {
				g.err = err
				return "an error accured: " + err.Error() + " press any keyboard key to exit"
			}

			if i == g.crrntRow && j == g.crrntCol {
				cols[j] = k.Render(g.currentKeyColor)
			} else {
				cols[j] = k.Render(g.keyColor)
			}
		}
		rows[i] = lipgloss.JoinHorizontal(lipgloss.Bottom, cols...)
	}
	questionList, err := data.GetGroupQuestions(g.usr.Group)
	table := lipgloss.JoinVertical(lipgloss.Center, rows...)
	questions := lipgloss.JoinVertical(lipgloss.Left, questionList...)
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
	if g.err != nil {
		return g, tea.Quit
	}
	if g.Ended() {
		if g.updateCounter < 1 {
			g.updateCounter++
			return g, nil
		} else {
      return g, g.gotoNextGame()
    }
	}
	switch msg := msg.(type) {
	case endGameMsg:
		//TODO show appropriate view end screen
		return g, g.gotoNextGame()
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

func (g *game) gotoNextGame() tea.Cmd {
	err := data.GroupGotoNextGame(g.usr.Group)
	if err != nil {
		g.err = err
		return func() tea.Msg {
			return errAccuredMsg{}
		}
	}

	var initialRow, initialCol int

	initialCol, err = data.GetGroupInitialCol(g.usr.Group)
	if err != nil {
		g.err = err
		return nil
	}

	initialRow, err = data.GetGroupInitialRow(g.usr.Group)
	if err != nil {
		g.err = err
		return nil
	}

	g.crrntCol = initialCol
	g.crrntRow = initialRow
	return nil
}

type errAccuredMsg struct{}

func (g *game) goDown() tea.Cmd {
	rowCount, err := data.GetGroupRows(g.usr.Group)
	if err != nil {
		g.err = err
		return func() tea.Msg {
			return errAccuredMsg{}
		}
	}
	if g.crrntRow+1 == rowCount {
		return nil
	}
	k, err := data.GetGroupRowColumn(g.usr.Group, g.crrntRow+1, g.crrntCol)
	if err != nil {
		g.err = err
		return func() tea.Msg {
			return errAccuredMsg{}
		}
	}
	if k.State == key.READONLY {
		return nil
	}
	g.crrntRow++

	return nil
}

func (g *game) goUp() tea.Cmd {
	if g.crrntRow == 0 {
		return nil
	}
	k, err := data.GetGroupRowColumn(g.usr.Group, g.crrntRow-1, g.crrntCol)
	if err != nil {
		g.err = err
		return func() tea.Msg {
			return errAccuredMsg{}
		}
	}
	if k.State == key.READONLY {
		return nil
	}
	g.crrntRow--

	return nil
}

func (g *game) goLeft() tea.Cmd {
	if g.crrntCol == 0 {
		return nil
	}
	k, err := data.GetGroupRowColumn(g.usr.Group, g.crrntRow, g.crrntCol-1)
	if err != nil {
		g.err = err
		return func() tea.Msg {
			return errAccuredMsg{}
		}
	}
	if k.State == key.READONLY {
		return nil
	}
	g.crrntCol--

	return nil
}

func (g *game) goRight() tea.Cmd {
	colCount, err := data.GetGroupCols(g.usr.Group)
	if err != nil {
		g.err = err
		return func() tea.Msg {
			return errAccuredMsg{}
		}
	}
	if g.crrntCol+1 == colCount {
		return nil
	}
	k, err := data.GetGroupRowColumn(g.usr.Group, g.crrntRow, g.crrntCol+1)
	if err != nil {
		g.err = err
		return func() tea.Msg {
			return errAccuredMsg{}
		}
	}
	if k.State == key.READONLY {
		return nil
	}
	g.crrntCol++

	return nil
}

func (g *game) insertKey(r rune) tea.Cmd {
	k, err := data.GetGroupRowColumn(g.usr.Group, g.crrntRow, g.crrntCol)
	if err != nil {
		g.err = err
		return func() tea.Msg {
			return errAccuredMsg{}
		}
	}

	if !unicode.IsLetter(r) || k.State == key.READONLY {
		return nil
	}

	r = unicode.ToUpper(r)
	char := key.Letters[r]
	k.Char = char
	err = data.GroupInsertKeyAt(g.usr.Group, k, g.crrntRow, g.crrntCol)

	if err != nil {
		g.err = err
		return func() tea.Msg {
			return errAccuredMsg{}
		}
	}

	if g.Ended() {
		g.updateCounter = 0
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
	ended, err := data.GroupGameEnded(g.usr.Group)

	if err != nil {
		g.err = err
	}

	return ended
}

func (g *game) EndGame() tea.Cmd {
	return func() tea.Msg {
		return endGameMsg{}
	}

}

func newGame(cfg config.Colors, height, width int, u user.User) (_ *game, err error) {
	var initialRow, initialCol int

	initialCol, err = data.GetGroupInitialCol(u.Group)
	if err != nil {
		return
	}

	initialRow, err = data.GetGroupInitialRow(u.Group)
	if err != nil {
		return
	}
	g := game{}
	g.height = height
	g.width = width
	g.questionBorderColor = cfg.QuestionBorderColor
	g.questionTextColor = cfg.QuestionTextColor
	g.currentKeyColor = cfg.TableSelectedKeyColor
	g.keyColor = cfg.TableEditableKeyColor
	g.passPhraseKeyColor = cfg.PassPhraseKeyColor
	g.usr = u
	g.crrntCol = initialCol
	g.crrntRow = initialRow
	return &g, nil
}

// TODO show time of other users realtime
// TODO cause everybody can register multiple times with different names and cheat \
// it would be nice if 1. start game at fixed time or 2. have multiple with same level of difficulity
// or limit registration, so only verified users will be able to play
