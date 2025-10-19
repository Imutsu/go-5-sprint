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
	personaldata.Personal
	Steps        int
	TrainingType string
	Duration     time.Duration
}

// Реализация метода Parse
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

	trainingType := strings.TrimSpace(parts[1])
	duration, err := time.ParseDuration(strings.TrimSpace(parts[2]))
	if err != nil {
		return fmt.Errorf("invalid duration format: %v", err)
	}
	if duration <= 0 {
		return fmt.Errorf("duration must be positive")
	}

	t.Steps = steps
	t.TrainingType = trainingType
	t.Duration = duration
	return nil
}

func (t Training) Print() {
	info, err := t.ActionInfo()
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	fmt.Print(info)
}

func (t Training) ActionInfo() (string, error) {
	if t.Weight <= 0 || t.Height <= 0 {
		return "", fmt.Errorf("personal data is required")
	}
	if t.Steps <= 0 || t.Duration <= 0 {
		return "", fmt.Errorf("invalid steps or duration")
	}

	distance := spentenergy.Distance(t.Steps, t.Height)
	speed := spentenergy.MeanSpeed(t.Steps, t.Height, t.Duration)

	var calories float64
	var err error

	switch strings.ToLower(t.TrainingType) {
	case "бег":
		calories, err = spentenergy.RunningSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
	case "ходьба":
		calories, err = spentenergy.WalkingSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
	default:
		return "", fmt.Errorf("неизвестный тип тренировки: %s", t.TrainingType)
	}

	if err != nil {
		return "", err
	}

	info := fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		t.TrainingType, t.Duration.Hours(), distance, speed, calories,
	)

	return info, nil
}
