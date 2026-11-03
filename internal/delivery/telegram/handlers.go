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
		return s.reply(c, formatError(err), mainMenu())
	}

	return s.reply(c, formatStartResponse(out), mainMenu())
}

func (s *Server) handleFlash(c tele.Context) error {
	return s.handleAscent(c, domain.OutcomeFlash)
}

func (s *Server) handleFail(c tele.Context) error {
	return s.handleAscent(c, domain.OutcomeFail)
}

func (s *Server) handleAscent(c tele.Context, outcome domain.Outcome) error {
	args := ""
	if c.Callback() == nil {
		args = strings.TrimSpace(c.Message().Payload)
	}
	if args == "" {
		return s.reply(c, formatGradePromptResponse(outcome), gradeMenu(outcome))
	}

	return s.saveAscent(c, outcome, domain.Grade(strings.ToUpper(args)))
}

// handleGrade handles a grade button press. Callback data: "{outcome}|{grade}".
func (s *Server) handleGrade(c tele.Context) error {
	args := c.Args()
	if len(args) != 2 {
		return s.reply(c, formatError(domain.ErrUnknownGrade), mainMenu())
	}

	outcome := domain.Outcome(args[0])
	if !outcome.Valid() {
		return s.reply(c, formatError(domain.ErrUnknownOutcome), mainMenu())
	}

	return s.saveAscent(c, outcome, domain.Grade(args[1]))
}

func (s *Server) saveAscent(c tele.Context, outcome domain.Outcome, grade domain.Grade) error {
	if !grade.Valid() {
		return s.reply(c, formatError(domain.ErrUnknownGrade), gradeMenu(outcome))
	}

	out, err := s.recordAscent.Exec(context.Background(), usecase.RecordAscentInput{
		TelegramID: c.Sender().ID,
		Grade:      grade,
		Outcome:    outcome,
	})
	if err != nil {
		return s.reply(c, formatError(err), mainMenu())
	}

	return s.reply(c, formatAscentResponse(out), mainMenu())
}

func (s *Server) handleHelp(c tele.Context) error {
	return s.reply(c, formatHelpResponse(), mainMenu())
}
