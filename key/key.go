package key

import (
	"errors"
	"strings"
	"unicode"

	"github.com/charmbracelet/lipgloss"
)

type key rune

func (k *key) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		return errors.New("len of key is 0")
	}
	var ok bool
	(*k), ok = Letters[unicode.ToUpper(rune(text[0]))]
	if !ok {
		return errors.New("invalid character " + string(text[0]))
	}
	return nil
}

const (
	A key = iota + 'A'
	B
	C
	D
	E
	F
	G
	H
	I
	J
	K
	L
	M
	N
	O
	P
	Q
	R
	S
	T
	U
	V
	W
	X
	Y
	Z
	EMPTY = ' '
)
const (
	ZERO key = iota + '0'
	ONE
	TWO
	THREE
	FOUR
	FIVE
	SIX
	SEVEN
	EIGHT
	NINE
)

type state int

func (s *state) UnmarshalText(text []byte) error {
	str := string(text)
	if strings.ToLower(str) == "r" {
		*s = READONLY
	} else if strings.ToLower(str) == "e" {
		*s = EDITABLE
	} else {
		return errors.New("state must be editable (e) or readonly (r), got " + str)
	}
	return nil
}

const (
	EDITABLE state = iota
	READONLY
)

type Key struct {
	Char  key   `json:"char"`
	State state `json:"state"`
}

func (k Key) Render(color lipgloss.Color) string {
	if k.State == READONLY {
		return lipgloss.NewStyle().
			Padding(0, 1).
			Border(lipgloss.HiddenBorder()).
			BorderForeground(lipgloss.NoColor{}).
			Foreground(lipgloss.NoColor{}).
			Render(string(k.Char))

	}
	return lipgloss.NewStyle().
		Padding(0, 1).
		Border(lipgloss.NormalBorder()).
		BorderForeground(color).
		Foreground(color).
		Render(string(k.Char))
}

func (k Key) IsEqual(j Key) bool {
	if k.Char == j.Char && k.State == j.State {
		return true
	}
	return false
}

var Letters = make(map[rune]key)

func init() {
	Letters['A'] = A
	Letters['B'] = B
	Letters['C'] = C
	Letters['D'] = D
	Letters['E'] = E
	Letters['F'] = F
	Letters['G'] = G
	Letters['H'] = H
	Letters['I'] = I
	Letters['J'] = J
	Letters['K'] = K
	Letters['L'] = L
	Letters['M'] = M
	Letters['N'] = N
	Letters['O'] = O
	Letters['P'] = P
	Letters['Q'] = Q
	Letters['R'] = R
	Letters['S'] = S
	Letters['T'] = T
	Letters['U'] = U
	Letters['V'] = V
	Letters['W'] = W
	Letters['X'] = X
	Letters['Y'] = Y
	Letters['Z'] = Z
	Letters[' '] = EMPTY
	Letters['0'] = ZERO
	Letters['1'] = ONE
	Letters['2'] = TWO
	Letters['3'] = THREE
	Letters['4'] = FOUR
	Letters['5'] = FIVE
	Letters['6'] = SIX
	Letters['7'] = SEVEN
	Letters['8'] = EIGHT
	Letters['9'] = NINE
}
