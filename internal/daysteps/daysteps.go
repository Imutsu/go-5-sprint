package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/Yandex-Practicum/tracker/internal/spentenergy"
)

type DaySteps struct {
	Steps    int
	Duration time.Duration
	Personal personaldata.Personal
	Weight   float64
	Height   float64
}

func (ds DaySteps) Print() {
	panic("unimplemented")
}

func (ds *DaySteps) Parse(dataString string) error {
	parts := strings.Split(dataString, ",")
	if len(parts) != 2 {
		return fmt.Errorf("invalid data format: expected 2 parts, got %d", len(parts))
	}

	// Строгая проверка пробелов - если есть пробелы внутри строки, это ошибка
	stepsStr := parts[0]
	if strings.TrimSpace(stepsStr) != stepsStr || strings.ContainsAny(stepsStr, " \t\n") {
		return fmt.Errorf("invalid steps format")
	}

	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return fmt.Errorf("invalid steps format: %v", err)
	}
	if steps <= 0 {
		return fmt.Errorf("steps must be positive")
	}

	// Строгая проверка пробелов для duration
	durationStr := parts[1]
	if strings.TrimSpace(durationStr) != durationStr || strings.ContainsAny(durationStr, " \t\n") {
		return fmt.Errorf("invalid duration format")
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return fmt.Errorf("invalid duration format: %v", err)
	}
	if duration <= 0 {
		return fmt.Errorf("duration must be positive")
	}

	ds.Steps = steps
	ds.Duration = duration
	return nil
}

func (ds DaySteps) ActionInfo() (string, error) {
	weight := ds.Personal.Weight
	height := ds.Personal.Height

	if ds.Weight > 0 {
		weight = ds.Weight
	}
	if ds.Height > 0 {
		height = ds.Height
	}

	if weight <= 0 || height <= 0 {
		return "", fmt.Errorf("personal data is required")
	}

	if ds.Steps <= 0 || ds.Duration <= 0 {
		return "", fmt.Errorf("invalid steps or duration")
	}

	distance := spentenergy.Distance(ds.Steps, height)
	calories, err := spentenergy.WalkingSpentCalories(
		ds.Steps,
		weight,
		height,
		ds.Duration,
	)
	if err != nil {
		return "", err
	}

	info := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		ds.Steps, distance, calories)

	return info, nil
}
