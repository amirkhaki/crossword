package data

import (
	"github.com/amirkhaki/crossword/key"
	"github.com/amirkhaki/crossword/user"
)

var d *Data

func init() {
	if d == nil {
		d = NewData()
	}
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
