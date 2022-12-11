package data

import (
	"fmt"
	"sort"
	"time"

	"github.com/amirkhaki/crossword/config"
	"github.com/amirkhaki/crossword/key"
	"github.com/amirkhaki/crossword/user"
)

type gameState struct {
	actual     [][]key.Key
	rows       int
	cols       int
	initialRow int
	initialCol int
	questions  []string
}

func (g gameState) ended() bool {
	for i := 0; i < g.rows; i++ {
		for j := 0; j < g.cols; j++ {
			if g.actual[i][j].Char != g.actual[i][j].MustBe {
				return false
			}
		}
	}
	return true
}

func (g gameState) isValidKey(row, col int) bool {
	if row >= g.rows {
		return false
	}

	if col >= g.cols {
		return false
	}

	return true
}

type GroupItem struct {
	startTime int64
	endTime   int64
	groupName string
}

func (g GroupItem) FilterValue() string {
	return g.groupName
}

func (g GroupItem) Title() string {
	return g.groupName
}

func (g GroupItem) Desciption() string {
	return fmt.Sprintf("Ended in %d seconds", g.endTime-g.startTime)
}

type Data struct {
	games map[user.Group]struct {
		states           []gameState
		currentGameIndex int
		isAfterGame      bool
		startTime        int64
		endTime          int64
		started          bool
		passphrase       string
	}
}

func (d *Data) GetItems() (l []GroupItem) {
	for k, v := range d.games {
		if v.endTime == 0 {
			continue
		}
		l = append(l, GroupItem{startTime: v.startTime, endTime: v.endTime, groupName: k.Name})
	}
	sort.Slice(l, func(i, j int) bool {
		return (l[i].endTime - l[i].startTime) <= (l[j].endTime - l[j].startTime)
	})
	return
}

func (d *Data) GroupIsPassphraseCorrect(grp user.Group, passphrase string) (_ bool, err error) {
	g, ok := d.games[grp]
	if !ok {
		err = GroupNotFoundError(fmt.Errorf("GetGroupInitialCol: Group not found"))
		return
	}
	if g.passphrase == passphrase {
		return true, nil
	}
	return false, nil
}

type GroupNotFoundError error

func (d *Data) GetGroupInitialCol(grp user.Group) (_ int, err error) {
	g, ok := d.games[grp]
	if !ok {
		err = GroupNotFoundError(fmt.Errorf("GetGroupInitialCol: Group not found"))
		return
	}

	return g.states[g.currentGameIndex].initialCol, nil
}

func (d *Data) GetGroupInitialRow(grp user.Group) (_ int, err error) {
	g, ok := d.games[grp]
	if !ok {
		err = GroupNotFoundError(fmt.Errorf("GetGroupInitialRow: Group not found"))
		return
	}

	return g.states[g.currentGameIndex].initialRow, nil
}

func (d *Data) GetGroupRows(grp user.Group) (_ int, err error) {
	g, ok := d.games[grp]
	if !ok {
		err = GroupNotFoundError(fmt.Errorf("GetGroupRows: Group not found"))
		return
	}

	return g.states[g.currentGameIndex].rows, nil
}

func (d *Data) GetGroupCols(grp user.Group) (_ int, err error) {
	g, ok := d.games[grp]
	if !ok {
		err = GroupNotFoundError(fmt.Errorf("GetGroupCols: Group not found"))
		return
	}

	return g.states[g.currentGameIndex].cols, nil
}

func (d *Data) GetGroupRowColumn(grp user.Group, row, col int) (k key.Key, err error) {
	g, ok := d.games[grp]
	if !ok {
		err = GroupNotFoundError(fmt.Errorf("GetGroupRowColumn: Group not found"))
		return
	}

	if !g.states[g.currentGameIndex].isValidKey(row, col) {
		err = fmt.Errorf("GetGroupRowColumn: invalid row col: %d, %d", row, col)
		return
	}

	return g.states[g.currentGameIndex].actual[row][col], nil
}

func (d *Data) GetGroupQuestions(grp user.Group) (_ []string, err error) {
	g, ok := d.games[grp]
	if !ok {
		err = GroupNotFoundError(fmt.Errorf("GetGroupQuestions: Group not found"))
		return
	}

	return g.states[g.currentGameIndex].questions, nil
}

func (d *Data) GroupInsertKeyAt(grp user.Group, k key.Key, row, col int) (err error) {
	g, ok := d.games[grp]
	if !ok {
		err = GroupNotFoundError(fmt.Errorf("GroupInsertKeyAt: Group not found"))
		return
	}

	if !g.states[g.currentGameIndex].isValidKey(row, col) {
		err = fmt.Errorf("GroupInsertKeyAt: invalid row col: %d, %d", row, col)
		return
	}

	g.states[g.currentGameIndex].actual[row][col] = k

	if !g.started {
		g.startTime = time.Now().Unix()
	}

	if g.states[g.currentGameIndex].ended() {
		g.isAfterGame = true
	}
	d.games[grp] = g
	return nil
}

func (d *Data) GroupIsAfterGame(grp user.Group) (_ bool, err error) {
	g, ok := d.games[grp]
	if !ok {
		err = GroupNotFoundError(fmt.Errorf("GroupIsAfterGame: Group not found"))
		return
	}

	return g.isAfterGame, nil
}

func (d *Data) GroupGameEnded(grp user.Group) (_ bool, err error) {
	g, ok := d.games[grp]
	if !ok {
		err = GroupNotFoundError(fmt.Errorf("GroupGameEnded: Group not found"))
		return
	}

	return g.states[g.currentGameIndex].ended(), nil
}

type GroupExistsError error

func (d *Data) AddGroup(grp user.Group, cfgs []config.Game, ps string) error {
	g, ok := d.games[grp]
	if ok {
		return GroupExistsError(fmt.Errorf("AddGroup: group already exists"))
	}
	for _, cfg := range cfgs {
		state := gameState{}
		state.questions = cfg.Questions
		state.rows = cfg.Rows
		state.cols = cfg.Cols
		state.initialCol = cfg.InitialCol
		state.initialRow = cfg.InitialRow
		state.actual = make([][]key.Key, cfg.Rows)
		for i := 0; i < cfg.Rows; i++ {
			state.actual[i] = make([]key.Key, cfg.Cols)
		}
		for _, k := range cfg.Actual.Keys {
			state.actual[k.Row][k.Col] = k.Key
		}
		g.states = append(g.states, state)
		g.passphrase = ps
	}
	d.games[grp] = g
	return nil
}

type AllGamesDoneError error

func (d *Data) GroupGotoNextGame(grp user.Group) error {
	g, ok := d.games[grp]
	if !ok {
		return GroupNotFoundError(fmt.Errorf("GroupGotoNextGame: Group not found"))
	}
	if len(g.states)-1 == g.currentGameIndex {
		g.endTime = time.Now().Unix()
		d.games[grp] = g
		return AllGamesDoneError(fmt.Errorf("GroupGotoNextGame: all games done"))
	}
	g.currentGameIndex++
	g.isAfterGame = false
	d.games[grp] = g
	return nil
}

func NewData() *Data {
	d := Data{}
	d.games = make(map[user.Group]struct {
		states           []gameState
		currentGameIndex int
		isAfterGame      bool
		startTime        int64
		endTime          int64
		started          bool
		passphrase       string
	})
	return &d
}
