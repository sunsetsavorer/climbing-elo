package telegram

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sunsetsavorer/climbing-elo/internal/domain"
	"github.com/sunsetsavorer/climbing-elo/internal/usecase"
)

func formatStartResponse(out *usecase.GetOrCreateClimberOutput) string {
	greeting := "С возвращением"
	if out.Created {
		greeting = "Привет"
	}
	return fmt.Sprintf("%s, @%s!\n\n%s", greeting, out.Climber.Nickname, formatProfile(out.Climber))
}

func formatAscentResponse(out *usecase.RecordAscentOutput) string {
	return fmt.Sprintf(
		"Пролаз записан!\n\nКатегория:  %s\nИзменение:  %s\nРейтинг:    %d",
		out.Ascent.Grade,
		formatDelta(out.Ascent.Delta),
		out.Ascent.RatingAfter,
	)
}

func formatAscentNoGradeResponse(outcome domain.Outcome) string {
	return fmt.Sprintf("Укажи категорию. Например: /%s 6B", string(outcome))
}

func formatHelpResponse() string {
	return "Команды:\n\n" +
		"/start — регистрация и профиль\n" +
		"/flash {категория} — флеш\n" +
		"/fail {категория} — не пролез\n" +
		"/stat — твой рейтинг\n" +
		"/help — список команд\n\n" +
		"Категории: " + formatGrades()
}

func formatStatResponse(out *usecase.GetOrCreateClimberOutput) string {
	return formatProfile(out.Climber)
}

func formatProfile(climber *domain.Climber) string {
	hardestGrade := "—"
	if climber.HardestGrade != nil {
		hardestGrade = string(*climber.HardestGrade)
	}
	return fmt.Sprintf(
		"Рейтинг:       %d\nПролазов:      %d\nЛучший флеш:   %s",
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
		return fmt.Sprintf("Неизвестная категория.\nДоступные: %s", formatGrades())
	case errors.Is(err, domain.ErrUnknownOutcome):
		return "Неизвестный результат. Используй /flash или /fail."
	case errors.Is(err, domain.ErrClimberNotFound):
		return "Скалолаз не найден. Начни с команды /start."
	default:
		return "Что-то пошло не так. Попробуй ещё раз."
	}
}

func formatGrades() string {
	grades := domain.AllGrades()
	names := make([]string, len(grades))
	for i, g := range grades {
		names[i] = string(g)
	}
	return strings.Join(names, ", ")
}
