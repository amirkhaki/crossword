package key

import "github.com/charmbracelet/lipgloss"

type key rune

func (k key) Render(color lipgloss.Color) string {
	if k == LOCKED{
		return lipgloss.NewStyle().
			Padding(0, 1).
			Border(lipgloss.HiddenBorder()).
			BorderForeground(lipgloss.NoColor{}).
			Foreground(lipgloss.NoColor{}).
			Render(string(' '))

	}
	return lipgloss.NewStyle().
		Padding(0, 1).
		Border(lipgloss.NormalBorder()).
		BorderForeground(color).
		Foreground(color).
		Render(string(k))
}

// TODO add support for locked letters, so table can have an initial info which users can't change it


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
	LOCKED
	EMPTY = ' '
)


const (
	ROWS = 13
	COLS = 13
)

type Table [ROWS][COLS]key


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
}
