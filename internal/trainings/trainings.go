package trainings

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/Yandex-Practicum/tracker/internal/spentenergy"
)

type Training struct {
	Steps        int
	TrainingType string
	Duration     time.Duration
	Personal     personaldata.Personal
	Weight       float64
	Height       float64
}

func (t Training) Print() {
	panic("unimplemented")
}

func (t *Training) Parse(dataString string) error {
	parts := strings.Split(dataString, ",")
	if len(parts) != 3 {
		return fmt.Errorf("invalid data format: expected 3 parts, got %d", len(parts))
	}

	steps, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return fmt.Errorf("invalid steps format: %v", err)
	}
	if steps <= 0 {
		return fmt.Errorf("steps must be positive")
	}

	// Парсим ЛЮБОЙ тип тренировки без проверки
	trainingType := strings.TrimSpace(parts[1])
	t.TrainingType = trainingType

	duration, err := time.ParseDuration(strings.TrimSpace(parts[2]))
	if err != nil {
		return fmt.Errorf("invalid duration format: %v", err)
	}
	if duration <= 0 {
		return fmt.Errorf("duration must be positive")
	}

	t.Steps = steps
	t.Duration = duration
	return nil
}

func (t Training) ActionInfo() (string, error) {
	weight := t.Personal.Weight
	height := t.Personal.Height

	if t.Weight > 0 {
		weight = t.Weight
	}
	if t.Height > 0 {
		height = t.Height
	}

	if weight <= 0 || height <= 0 {
		return "", fmt.Errorf("personal data is required")
	}

	if t.Steps <= 0 || t.Duration <= 0 {
		return "", fmt.Errorf("invalid steps or duration")
	}

	distance := spentenergy.Distance(t.Steps, height)
	speed := spentenergy.MeanSpeed(t.Steps, height, t.Duration)

	var calories float64
	var err error

	// Проверяем тип тренировки только в ActionInfo, а не в Parse
	switch strings.ToLower(t.TrainingType) {
	case "бег", "run", "running":
		calories, err = spentenergy.RunningSpentCalories(
			t.Steps,
			weight,
			height,
			t.Duration,
		)
	case "ходьба", "walk", "walking":
		calories, err = spentenergy.WalkingSpentCalories(
			t.Steps,
			weight,
			height,
			t.Duration,
		)
	default:
		return "", fmt.Errorf("неизвестный тип тренировки: %s", t.TrainingType)
	}

	if err != nil {
		return "", err
	}

	info := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		t.TrainingType, t.Duration.Hours(), distance, speed, calories)

	return info, nil
}
