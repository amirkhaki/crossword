package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/amirkhaki/crossword/config"
	"github.com/amirkhaki/crossword/model"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	lm "github.com/charmbracelet/wish/logging"
	"github.com/gliderlabs/ssh"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var withServer *bool
var configPath *string
var serverHost *string
var serverPort *int

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, active := s.Pty()
	if !active {
		wish.Fatalln(s, "no active terminal, skipping")
		return nil, nil
	}
	cfg, err := config.New(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	g, err := model.NewGame(cfg, pty.Window.Height, pty.Window.Width)
	if err != nil {
		log.Fatal(err)
	}
	return g, []tea.ProgramOption{tea.WithAltScreen()}
}

func init() {
	configPath = flag.String("config", "config.json", "path to config file, format must be json")
	withServer = flag.Bool("server", false, "whether run ssh server or not")
	serverHost = flag.String("host", "127.0.0.1", "host for server")
	serverPort = flag.Int("port", 2222, "port for server")
}

// having a map of [user][]program
// any interaction in a game will be sent to all of programs(p.Send)
// so state of game should not be saved in the model cause it is user specific
// thinking about having a map[user]struct {[]program, state}

func main() {
	flag.Parse()
	cfg, err := config.New(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	g, err := model.NewGame(cfg, 0, 0)
	if err != nil {
		log.Fatal(err)
	}
	if *withServer {
		s, err := wish.NewServer(
			wish.WithAddress(fmt.Sprintf("%s:%d", *serverHost, *serverPort)),
			wish.WithHostKeyPath(".ssh/term_info_ed25519"),
			wish.WithMiddleware(
				bm.Middleware(teaHandler),
				lm.Middleware(),
			),
		)
		if err != nil {
			log.Fatalln(err)
		}
		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		log.Printf("Starting SSH server on %s:%d", *serverHost, *serverPort)
		go func() {
			if err = s.ListenAndServe(); err != nil {
				log.Fatalln(err)
			}
		}()

		<-done
		log.Println("Stopping SSH server")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer func() { cancel() }()
		if err := s.Shutdown(ctx); err != nil {
			log.Fatalln(err)
		}

	} else {
		p := tea.NewProgram(g)
		if err := p.Start(); err != nil {
			log.Fatal(err)
		}
	}
}
