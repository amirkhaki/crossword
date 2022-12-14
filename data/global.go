package data

import (
	"github.com/amirkhaki/crossword/config"
	"github.com/amirkhaki/crossword/key"
	"github.com/amirkhaki/crossword/user"
)

var d *Data

func init() {
	if d == nil {
		d = NewData()
	}
}

func GroupAllGameEnded(grp user.Group) (bool, error) {
	return d.GroupAllGameEnded(grp)
}

func GroupEndAllGame(grp user.Group) error {
	return d.GroupEndAllGame(grp)
}

func GroupIsPassphraseCorrect(grp user.Group, passphrase string) (bool, error) {
	return d.GroupIsPassphraseCorrect(grp, passphrase)
}

func GetItems() []GroupItem {
	return d.GetItems()
}

func IsAfterGame(grp user.Group) (bool, error) {
	return d.GroupIsAfterGame(grp)
}

func GetGroupRows(grp user.Group) (int, error) {
	return d.GetGroupRows(grp)
}
func GetGroupCols(grp user.Group) (int, error) {
	return d.GetGroupCols(grp)
}

func GetGroupRowColumn(grp user.Group, row, col int) (k key.Key, err error) {
	return d.GetGroupRowColumn(grp, row, col)
}

func GetGroupQuestions(grp user.Group) ([]string, error) {
	return d.GetGroupQuestions(grp)
}
func GroupInsertKeyAt(grp user.Group, k key.Key, row, col int) (err error) {
	return d.GroupInsertKeyAt(grp, k, row, col)
}

func AddGroup(grp user.Group, cfgs []config.Game, ps string) error {
	return d.AddGroup(grp, cfgs, ps)
}

func GroupGotoNextGame(grp user.Group) error {
	return d.GroupGotoNextGame(grp)
}

func GetGroupInitialRow(grp user.Group) (int, error) {
	return d.GetGroupInitialRow(grp)
}

func GetGroupInitialCol(grp user.Group) (int, error) {
	return d.GetGroupInitialCol(grp)
}
func GroupGameEnded(grp user.Group) (bool, error) {
	return d.GroupGameEnded(grp)
}
