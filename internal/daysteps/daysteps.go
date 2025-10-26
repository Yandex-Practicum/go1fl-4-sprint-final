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
	dataToStrings := strings.Split(data, ",")
	if len(dataToStrings) != 2 {
		return 0, 0, errors.New("not enough data to parse a package")
	}
	stepsRaw := dataToStrings[0]
	if stepsRaw != strings.TrimSpace(stepsRaw) {
        return 0, 0, errors.New("invalid steps format")
    }
	steps, err := strconv.Atoi(stepsRaw)
	if err != nil {
		return 0, 0, errors.New("failed to convert steps to int: " + err.Error())
	}
	if steps < 1 {
		return 0, 0, errors.New("can't execute program if there are less than 1 steps")
	}

	duration, err := time.ParseDuration(strings.TrimSpace(dataToStrings[1]))
	if err != nil {
		return 0, 0, errors.New("failed to parse duration: " + err.Error())
	}
	if duration <= 0 {
        return 0, 0, errors.New("duration must be > 0")
    }
	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) (string, error) {
    steps, duration, err := parsePackage(data)
    if err != nil {
        log.Printf("failed to retrieve data: %v", err)
        return "", nil
    }

    distance := (float64(steps) * stepLength) / mInKm

    calBurned, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
    if err != nil {
        log.Printf("failed to calculate spent calories: %v", err)
        return "", nil
    }

    report := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distance, calBurned)
    return report, nil
}