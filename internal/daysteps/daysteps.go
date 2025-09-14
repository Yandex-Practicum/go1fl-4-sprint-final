package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	// TODO: реализовать функцию
	parts := strings.Split(strings.TrimSpace(data), ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid format: expected '<steps>,<training duration>'")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid steps data %q: %w", parts[0], err)
	}

	if steps < 0 {
		return 0, 0, fmt.Errorf("invalid steps data %q", parts[0])
	}

	durStr := strings.TrimSpace(parts[1])
	if durStr == "" {
		return 0, 0, fmt.Errorf("duration part is empty")
	}

	dur, err := time.ParseDuration(durStr)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid duration %q: %w", durStr, err)
	}
	return steps, dur, nil
}

func DayActionInfo(data string, weight, height float64) string {
	// TODO: реализовать функцию
	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println("Failed to parse package:", err)
		return ""
	}

	distance := float64(steps) * stepLength
	distanceKm := distance / float64(mInKm)
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)

	result := fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.",
		steps,
		distanceKm,
		calories,
	)

	return result
}
