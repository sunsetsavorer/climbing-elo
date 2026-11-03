package telegram

import (
	"fmt"
	"sync"

	tele "gopkg.in/telebot.v3"
)

type Server struct {
	bot                *tele.Bot
	getOrCreateClimber GetOrCreateClimberUseCase
	recordAscent       RecordAscentUseCase
	lastReplies        sync.Map // chat ID -> tele.StoredMessage
}

func NewServer(
	bot *tele.Bot,
	getOrCreateClimber GetOrCreateClimberUseCase,
	recordAscent RecordAscentUseCase,
) (*Server, error) {
	if bot == nil {
		return nil, ErrBotRequired
	}
	if getOrCreateClimber == nil {
		return nil, ErrGetOrCreateClimberRequired
	}
	if recordAscent == nil {
		return nil, ErrRecordAscentRequired
	}
	s := &Server{
		bot:                bot,
		getOrCreateClimber: getOrCreateClimber,
		recordAscent:       recordAscent,
	}
	s.registerHandlers()
	if err := s.bot.SetCommands(botCommands()); err != nil {
		return nil, fmt.Errorf("set bot commands: %w", err)
	}
	return s, nil
}

type command struct {
	name        string
	args        string // usage hint for /help, e.g. "категория"
	description string
}

// commands are shown in the Telegram command menu, "/" suggestions and /help.
var commands = []command{
	{name: "start", description: "Регистрация и профиль"},
	{name: "flash", args: "категория", description: "Флеш"},
	{name: "fail", args: "категория", description: "Не пролез"},
	{name: "stat", description: "Твой рейтинг"},
	{name: "help", description: "Список команд"},
}

func botCommands() []tele.Command {
	out := make([]tele.Command, len(commands))
	for i, cmd := range commands {
		out[i] = tele.Command{Text: cmd.name, Description: cmd.description}
	}
	return out
}

func (s *Server) registerHandlers() {
	s.bot.Handle("/start", s.handleStart)
	s.bot.Handle("/help", s.handleHelp)
	s.bot.Handle("/flash", s.handleFlash)
	s.bot.Handle("/fail", s.handleFail)
	s.bot.Handle("/stat", s.handleStart)

	s.bot.Handle(&btnFlash, s.handleFlash)
	s.bot.Handle(&btnFail, s.handleFail)
	s.bot.Handle(&btnStat, s.handleStart)
	s.bot.Handle(&btnMenu, s.handleStart)

	s.bot.Handle(&btnGrade, s.handleGrade)
}

func (s *Server) Start() {
	s.bot.Start()
}

func (s *Server) Stop() {
	s.bot.Stop()
}
