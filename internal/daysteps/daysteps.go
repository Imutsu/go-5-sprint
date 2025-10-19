package daysteps

import (
	"errors"
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
	personaldata.Personal // Встраивание Personal
	Weight   float64
	Height   float64
}

func (ds DaySteps) Print() {
	fmt.Printf("Шаги: %d, Длительность: %v\n", ds.Steps, ds.Duration)
	if ds.Weight > 0 {
		fmt.Printf("Вес: %.2f\n", ds.Weight)
	}
	if ds.Height > 0 {
		fmt.Printf("Рост: %.2f\n", ds.Height)
	}
	// Если у Personal есть метод Print, иначе можно удалить эту строку
	// ds.Personal.Print()
}

func (ds *DaySteps) Parse(dataString string) error {
	parts := strings.Split(dataString, ",")
	if len(parts) != 2 {
		return fmt.Errorf("invalid data format: expected 2 parts, got %d", len(parts))
	}

	rawSteps := parts[0]

	if rawSteps != strings.TrimSpace(rawSteps) {
		return fmt.Errorf("invalid steps format: contains leading/trailing spaces")
	}

	steps, err := strconv.Atoi(rawSteps)
	if err != nil {
		return fmt.Errorf("invalid steps format: %w", err)
	}
	if steps <= 0 {
		return errors.New("steps must be positive")
	}

	durationStr := strings.TrimSpace(parts[1])
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return fmt.Errorf("invalid duration format: %w", err)
	}
	if duration <= 0 {
		return errors.New("duration must be positive")
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
		return "", errors.New("personal data is required") // Используем errors.New когда не нужно форматирование
	}

	if ds.Steps <= 0 || ds.Duration <= 0 {
		return "", errors.New("invalid steps or duration") // Используем errors.New когда не нужно форматирование
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