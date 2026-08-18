package daysteps

import (
	"errors"
	"fmt"
	"log"
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
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, errors.New("invalid format: expected 2 parts")
	}

	stepsStr := parts[0]
	durationStr := parts[1]

	if stepsStr != strings.TrimSpace(stepsStr) || durationStr != strings.TrimSpace(durationStr) {
		return 0, 0, errors.New("invalid format: whitespace detected")
	}

	if strings.Contains(stepsStr, " ") {
		return 0, 0, errors.New("invalid steps format")
	}

	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, 0, err
	}

	if steps <= 0 {
		return 0, 0, errors.New("steps must be positive")
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, 0, err
	}

	if duration <= 0 {
		return 0, 0, errors.New("duration must be positive")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}

	distKm := float64(steps) * stepLength / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println(err)
		return ""
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps,
		distKm,
		calories,
	)
}
