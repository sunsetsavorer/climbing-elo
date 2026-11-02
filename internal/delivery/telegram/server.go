package telegram

import (
	tele "gopkg.in/telebot.v3"
)

type Server struct {
	bot                *tele.Bot
	getOrCreateClimber GetOrCreateClimberUseCase
	recordAscent       RecordAscentUseCase
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
	return s, nil
}

func (s *Server) registerHandlers() {
	s.bot.Handle("/start", s.handleStart)
	s.bot.Handle("/help", s.handleHelp)
	s.bot.Handle("/flash", s.handleFlash)
	s.bot.Handle("/fail", s.handleFail)
	s.bot.Handle("/stat", s.handleStat)
}

func (s *Server) Start() {
	s.bot.Start()
}

func (s *Server) Stop() {
	s.bot.Stop()
}
