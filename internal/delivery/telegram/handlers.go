package telegram

import (
	"context"
	"strings"

	"github.com/sunsetsavorer/climbing-elo/internal/domain"
	"github.com/sunsetsavorer/climbing-elo/internal/usecase"
	tele "gopkg.in/telebot.v3"
)

func (s *Server) handleStart(c tele.Context) error {
	out, err := s.getOrCreateClimber.Exec(context.Background(), usecase.GetOrCreateClimberInput{
		TelegramID: c.Sender().ID,
		Nickname:   c.Sender().Username,
	})
	if err != nil {
		return c.Send(formatError(err))
	}

	return c.Send(formatStartResponse(out))
}

func (s *Server) handleFlash(c tele.Context) error {
	return s.handleAscent(c, domain.OutcomeFlash)
}

func (s *Server) handleFail(c tele.Context) error {
	return s.handleAscent(c, domain.OutcomeFail)
}

func (s *Server) handleAscent(c tele.Context, outcome domain.Outcome) error {
	args := strings.TrimSpace(c.Message().Payload)
	if args == "" {
		return c.Send(formatAscentNoGradeResponse(outcome))
	}

	grade := domain.Grade(strings.ToUpper(args))
	if !grade.Valid() {
		return c.Send(formatError(domain.ErrUnknownGrade))
	}

	out, err := s.recordAscent.Exec(context.Background(), usecase.RecordAscentInput{
		TelegramID: c.Sender().ID,
		Grade:      grade,
		Outcome:    outcome,
	})
	if err != nil {
		return c.Send(formatError(err))
	}

	return c.Send(formatAscentResponse(out))
}

func (s *Server) handleHelp(c tele.Context) error {
	return c.Send(formatHelpResponse())
}

func (s *Server) handleStat(c tele.Context) error {
	out, err := s.getOrCreateClimber.Exec(context.Background(), usecase.GetOrCreateClimberInput{
		TelegramID: c.Sender().ID,
		Nickname:   c.Sender().Username,
	})
	if err != nil {
		return c.Send(formatError(err))
	}

	return c.Send(formatStatResponse(out))
}
