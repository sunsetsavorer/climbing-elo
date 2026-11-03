package telegram

import (
	"github.com/sunsetsavorer/climbing-elo/internal/domain"
	tele "gopkg.in/telebot.v3"
)

const gradeButtonsPerRow = 3

// Endpoints for callback handlers. Text is set per markup.
var (
	btnFlash = tele.Btn{Unique: "flash", Text: "⚡ Флеш"}
	btnFail  = tele.Btn{Unique: "fail", Text: "✖ Не пролез"}
	btnStat  = tele.Btn{Unique: "stat", Text: "📊 Статистика"}
	btnMenu  = tele.Btn{Unique: "menu", Text: "← Назад"}
	btnGrade = tele.Btn{Unique: "grade"}
)

func mainMenu() *tele.ReplyMarkup {
	m := &tele.ReplyMarkup{}
	m.Inline(
		m.Row(btnFlash, btnFail),
		m.Row(btnStat),
	)
	return m
}

func gradeMenu(outcome domain.Outcome) *tele.ReplyMarkup {
	m := &tele.ReplyMarkup{}
	grades := domain.AllGrades()
	btns := make([]tele.Btn, len(grades))
	for i, g := range grades {
		btns[i] = m.Data(string(g), btnGrade.Unique, string(outcome), string(g))
	}
	rows := m.Split(gradeButtonsPerRow, btns)
	rows = append(rows, m.Row(btnMenu))
	m.Inline(rows...)
	return m
}
