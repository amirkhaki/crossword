package model

import (
	"fmt"
	"log"
	"time"
	"unicode"

	"github.com/amirkhaki/crossword/config"
	"github.com/amirkhaki/crossword/data"
	"github.com/amirkhaki/crossword/key"
	"github.com/amirkhaki/crossword/user"

	"github.com/charmbracelet/bubbles/textinput"
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
	height int
	width  int
	inited bool
}

func (_ endScreen) Init() tea.Cmd {
	return nil
}

type tickMsg struct{}

func doTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (e endScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			cmd = tea.Quit
		}
	case tickMsg:
		cmd = doTick()
	default:
		if !e.inited {
			e.inited = true
			cmd = doTick()
		}
	}
	return e, cmd
}

func (e endScreen) View() string {
	style := lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62"))
	var rows []string
	for _, v := range data.GetItems() {
		rows = append(rows, style.Render(fmt.Sprintf("%s\n%s", v.Title(), v.Desciption())))
	}
	return lipgloss.Place(e.width, e.height, lipgloss.Center, lipgloss.Center,
		style.Render(lipgloss.JoinVertical(lipgloss.Center, rows...)))

}

type passphraseScreen struct {
	width      int
	height     int
	usr        user.User
	passphrase textinput.Model
}

func (ps passphraseScreen) Init() tea.Cmd {
	return nil
}

func (ps passphraseScreen) doResize(msg tea.WindowSizeMsg) passphraseScreen {
	ps.height = msg.Height
	ps.width = msg.Width
	return ps
}
func (ps passphraseScreen) checkAnswer() (tea.Model, tea.Cmd) {
	ok, err := data.GroupIsPassphraseCorrect(ps.usr.Group, ps.passphrase.Value())
	if err != nil || !ok {
		log.Println(err)
		return ps, nil
	}
	err = data.GroupEndAllGame(ps.usr.Group)
	if err != nil {
		// TODO handle error and show it to user
		log.Println(err)
		return ps, nil
	}
	log.Println("checkAnswer game ended")
	return endScreen{height: ps.height, width: ps.width}.Update(nil)

}

func (ps passphraseScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	ok, err := data.GroupAllGameEnded(ps.usr.Group)
	if err == nil && ok {
		return endScreen{height: ps.height, width: ps.width}.Update(nil)
	}
	ps.passphrase.Focus()
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl-c":
			return ps, tea.Quit
		case "enter":
			return ps.checkAnswer()
		}
	case tea.WindowSizeMsg:
		ps = ps.doResize(msg)
		return ps, nil
	}
	var cmd tea.Cmd
	ps.passphrase, cmd = ps.passphrase.Update(msg)
	return ps, cmd
}

func (ps passphraseScreen) View() string {
	return lipgloss.Place(ps.width, ps.height, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Left, "Please enter passphrase", ps.passphrase.View()))
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
	if _, ok := msg.(AllDoneMsg); ok {
		mdl := textinput.New()
		return passphraseScreen{height: g.height, width: g.width, usr: g.usr, passphrase: mdl}.Update(nil)
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

type AllDoneMsg struct{}

func (g *game) gotoNextGame() tea.Cmd {
	err := data.GroupGotoNextGame(g.usr.Group)
	if err != nil {
		if _, ok := err.(data.AllGamesDoneError); ok {
			return func() tea.Msg {
				return AllDoneMsg{}
			}
		}
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
