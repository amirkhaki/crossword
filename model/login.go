package model

import (
	"context"

	"github.com/amirkhaki/crossword/config"
	"github.com/amirkhaki/crossword/storage"
	"github.com/amirkhaki/crossword/user"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type login struct {
	cfg      config.Config
	status   string
	height   int
	width    int
	username textinput.Model
	password textinput.Model
}

func (l login) Init() tea.Cmd {
	return nil
}

func (l login) loginUser() (tea.Model, tea.Cmd) {
	username := l.username.Value()
	password := l.password.Value()
	u, err := storage.Store.GetUser(context.Background(),
		user.NewUser(username, password, user.Group{}), func(u1, u2 user.User) bool {
			if u1.Username == u2.Username && u1.Password == u2.Password {
				return true
			}
			return false
		})
	if err != nil {
		_, ok := err.(storage.UserNotFoundError)
		form := NewLogin(l.cfg, l.height, l.width)
		if ok {
			form.status = "invalid username and/or password! try again"
		} else {
			form.status = "an error accured: " + err.Error()
		}
		return form, nil
	}
	g, err := newGame(l.cfg.Colors, l.height, l.width, u)
	if err != nil {
		form := NewLogin(l.cfg, l.height, l.width)
		form.status = "an error accured: " + err.Error()
		return form, nil
	}
	return g.Update(nil)
}

func (l login) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if l.username.Focused() {
				l.username.Blur()
				l.password.Focus()
				return l, textinput.Blink
			} else {
				return l.loginUser()
			}
		}
	case tea.WindowSizeMsg:
		return l, l.doResize(msg)
	}
	var cmd tea.Cmd

	if l.username.Focused() {
		l.username, cmd = l.username.Update(msg)
	} else {
		l.password, cmd = l.password.Update(msg)
	}

	return l, cmd
}

func (g login) doResize(msg tea.WindowSizeMsg) tea.Cmd {
	g.height = msg.Height
	g.width = msg.Width
	return nil
}
func (l login) View() string {
	var status string
	if l.status != "" {
		status = lipgloss.NewStyle().Margin(1).Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#ff0000")).Render(l.status)
	}
	return lipgloss.Place(l.width, l.height, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Left, status, l.username.View(), l.password.View()))
}

func NewLogin(cfg config.Config, height, width int) login {
	l := login{height: height, width: width}
	l.username = textinput.New()
	l.username.Placeholder = "username"
	l.username.Focus()
	l.password = textinput.New()
	l.password.Placeholder = "password"
	l.cfg = cfg
	return l
}
