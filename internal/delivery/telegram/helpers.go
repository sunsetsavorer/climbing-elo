package telegram

import (
	"errors"
	"fmt"
	"html"
	"strings"

	"github.com/sunsetsavorer/climbing-elo/internal/domain"
	"github.com/sunsetsavorer/climbing-elo/internal/usecase"
)

// All responses are sent with HTML parse mode: escape user-provided text.

func formatStartResponse(out *usecase.GetOrCreateClimberOutput) string {
	greeting := "👋 <b>С возвращением"
	if out.Created {
		greeting = "👋 <b>Привет"
	}
	if out.Climber.Nickname != "" {
		greeting += ", @" + html.EscapeString(out.Climber.Nickname)
	}
	greeting += "!</b>"

	if out.Created {
		greeting += "\n\nЗаписывай пролазы, а я посчитаю твой рейтинг."
	}
	return greeting + "\n\n" + formatProfile(out.Climber)
}

func formatAscentResponse(out *usecase.RecordAscentOutput) string {
	a := out.Ascent

	title := fmt.Sprintf("⚡ <b>Флеш %s записан!</b>", a.Grade)
	if a.Outcome == domain.OutcomeFail {
		title = fmt.Sprintf("✖ <b>Попытка %s записана</b>", a.Grade)
	}

	trend := "📈"
	if a.Delta < 0 {
		trend = "📉"
	}

	return fmt.Sprintf(
		"%s\n\n%s Рейтинг: %d → <b>%d</b> (%s)",
		title,
		trend,
		a.RatingAfter-a.Delta,
		a.RatingAfter,
		formatDelta(a.Delta),
	)
}

func formatGradePromptResponse(outcome domain.Outcome) string {
	title := "⚡ <b>Флеш</b>"
	if outcome == domain.OutcomeFail {
		title = "✖ <b>Не пролез</b>"
	}
	return title + "\n\nВыбери категорию:"
}

func formatHelpResponse() string {
	var b strings.Builder
	b.WriteString("📖 <b>Как пользоваться</b>\n\n")
	b.WriteString("Жми кнопки под сообщением или пиши команды:\n\n")
	for _, cmd := range commands {
		b.WriteString("/" + cmd.name)
		if cmd.args != "" {
			b.WriteString(" <i>" + cmd.args + "</i>")
		}
		b.WriteString(" — " + strings.ToLower(cmd.description) + "\n")
	}
	b.WriteString("\n<b>Категории:</b>\n" + formatGrades())
	return b.String()
}

var leaderboardMedals = []string{"🥇", "🥈", "🥉"}

// formatLeaderboardResponse highlights the row of the climber with telegramID.
func formatLeaderboardResponse(out *usecase.GetLeaderboardOutput, telegramID int64) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("🏆 <b>Лидерборд</b>\n\nПервые %d по рейтингу:\n\n", leaderboardSize))

	if len(out.Climbers) == 0 {
		b.WriteString("Пока никого нет. Запиши первый пролаз!")
		return b.String()
	}

	for i, climber := range out.Climbers {
		place := fmt.Sprintf("%d.", i+1)
		if i < len(leaderboardMedals) {
			place = leaderboardMedals[i]
		}

		row := fmt.Sprintf("%s %s — %d", place, formatNickname(climber), climber.Rating)
		if climber.TelegramID == telegramID {
			row = "<b>" + row + "</b>"
		}
		b.WriteString(row + "\n")
	}
	return b.String()
}

func formatNickname(climber *domain.Climber) string {
	if climber.Nickname == "" {
		return "без ника"
	}
	return "@" + html.EscapeString(climber.Nickname)
}

func formatProfile(climber *domain.Climber) string {
	hardestGrade := "—"
	if climber.HardestGrade != nil {
		hardestGrade = string(*climber.HardestGrade)
	}
	return fmt.Sprintf(
		"⭐ Рейтинг: <b>%d</b>\n🧗 Пролазов: <b>%d</b>\n💪 Сложнейший грейд: <b>%s</b>",
		climber.Rating,
		climber.AscentsCount,
		hardestGrade,
	)
}

func formatDelta(delta int) string {
	if delta > 0 {
		return fmt.Sprintf("+%d", delta)
	}
	return fmt.Sprintf("%d", delta)
}

func formatError(err error) string {
	switch {
	case errors.Is(err, domain.ErrUnknownGrade):
		return "⚠️ <b>Неизвестная категория</b>\n\nДоступные:\n" + formatGrades()
	case errors.Is(err, domain.ErrUnknownOutcome):
		return "⚠️ <b>Неизвестный результат</b>\n\nИспользуй /flash или /fail."
	case errors.Is(err, domain.ErrClimberNotFound):
		return "⚠️ <b>Скалолаз не найден</b>\n\nНачни с команды /start."
	default:
		return "⚠️ <b>Что-то пошло не так</b>\n\nПопробуй ещё раз."
	}
}

func formatGrades() string {
	grades := domain.AllGrades()
	names := make([]string, len(grades))
	for i, g := range grades {
		names[i] = "<code>" + string(g) + "</code>"
	}
	return strings.Join(names, " ")
}
