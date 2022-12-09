package data

import (
	"fmt"

	"github.com/amirkhaki/crossword/key"
	"github.com/amirkhaki/crossword/user"
)

type gameState struct {
	actual    [][]key.Key
	rows      int
	cols      int
	questions []string
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

type Data struct {
	games map[user.Group]struct {
		states           []gameState
		currentGameIndex int
		isAfterGame      bool
	}
}

type GroupNotFoundError error

func (d *Data) GetGroupRows(grp user.Group) (_ int, err error) {
	g, ok := d.games[grp]
	if !ok {
		err = GroupNotFoundError(fmt.Errorf("Group not found"))
		return
	}

	return g.states[g.currentGameIndex].rows, nil
}

func (d *Data) GetGroupCols(grp user.Group) (_ int, err error) {
	g, ok := d.games[grp]
	if !ok {
		err = GroupNotFoundError(fmt.Errorf("Group not found"))
		return
	}

	return g.states[g.currentGameIndex].cols, nil
}

func (d *Data) GetGroupRowColumn(grp user.Group, row, col int) (k key.Key, err error) {
	g, ok := d.games[grp]
	if !ok {
		err = GroupNotFoundError(fmt.Errorf("Group not found"))
		return
	}

	if !g.states[g.currentGameIndex].isValidKey(row, col) {
		err = fmt.Errorf("data: invalid row col: %d, %d", row, col)
		return
	}

	return g.states[g.currentGameIndex].actual[row][col], nil
}

func (d *Data) GetGroupQuestions(grp user.Group) (_ []string, err error) {
	g, ok := d.games[grp]
	if !ok {
		err = GroupNotFoundError(fmt.Errorf("Group not found"))
		return
	}

	return g.states[g.currentGameIndex].questions, nil
}

func (d *Data) GroupInsertKeyAt(grp user.Group, k key.Key, row, col int) (err error) {
	g, ok := d.games[grp]
	if !ok {
		err = GroupNotFoundError(fmt.Errorf("Group not found"))
		return
	}

	if !g.states[g.currentGameIndex].isValidKey(row, col) {
		err = fmt.Errorf("data: invalid row col: %d, %d", row, col)
		return
	}

	g.states[g.currentGameIndex].actual[row][col] = k
	d.games[grp] = g
	return nil

}

func (d *Data) GroupGameEnded(grp user.Group) (_ bool, err error) {
	g, ok := d.games[grp]
	if !ok {
		err = GroupNotFoundError(fmt.Errorf("Group not found"))
		return
	}

  return g.states[g.currentGameIndex].ended(), nil
}

//TODO insert key and initialise data

func NewData() *Data {
	return &Data{}
}
